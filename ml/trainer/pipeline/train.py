from datetime import datetime, timezone

try:
    from trainer.config.config import Config
    from trainer.data.reader import load_latest_dataset
    from trainer.data.splitter import split_train_validation
    from trainer.evaluation.metrics import evaluate_binary_classifier
    from trainer.export.writer import export_model_bundle
    from trainer.features.encoder import fit_feature_encoder
    from trainer.models.catalog import build_candidate_models
    from trainer.registry.writer import update_registry
    from trainer.run_context import TrainerRunContext
except ImportError:
    from config.config import Config
    from data.reader import load_latest_dataset
    from data.splitter import split_train_validation
    from evaluation.metrics import evaluate_binary_classifier
    from export.writer import export_model_bundle
    from features.encoder import fit_feature_encoder
    from models.catalog import build_candidate_models
    from registry.writer import update_registry
    from run_context import TrainerRunContext

def run_training(config: Config, run_context: TrainerRunContext) -> None:
    print("loading transformed dataset")
    dataset, dataset_path = load_latest_dataset(
        run_context.datasets_dir,
        config.dataset.dataset_glob,
    )
    print(f"dataset path: {dataset_path}")
    print(f"dataset rows before filtering: {len(dataset)}")

    if "team_A_won" not in dataset.columns:
        raise ValueError("transformed dataset is missing required column: team_A_won")

    if dataset.empty:
        raise ValueError("transformed dataset is empty for the selected run window")

    if config.dataset.drop_draws:
        dataset = dataset[dataset["team_A_won"] != 0.5].copy()
        print(f"dataset rows after dropping draws: {len(dataset)}")

    if dataset.empty:
        raise ValueError("no usable training rows remain after filtering")

    train_frame, validation_frame = split_train_validation(dataset, config.split.validation_day_pct)
    print(f"train rows: {len(train_frame)}")
    print(f"validation rows: {len(validation_frame)}")

    encoder = fit_feature_encoder(train_frame)
    x_train = encoder.transform(train_frame)
    x_validation = encoder.transform(validation_frame)
    y_train = train_frame["team_A_won"].astype(int)
    y_validation = validation_frame["team_A_won"].astype(int)

    feature_names = encoder.feature_names()
    print(f"feature count: {len(feature_names)}")

    candidates = {}
    for model_name, model in build_candidate_models(config):
        print(f"training model: {model_name}")
        model.fit(x_train, y_train)
        validation_proba = model.predict_proba(x_validation)[:, 1]
        metrics = evaluate_binary_classifier(y_validation, validation_proba)
        print(f"{model_name} validation metrics: {metrics}")
        candidates[model_name] = {
            "model": model,
            "metrics": metrics,
        }

    selected_model_name = min(
        candidates.items(),
        key=lambda item: item[1]["metrics"]["log_loss"],
    )[0]
    selected = candidates[selected_model_name]
    print(f"selected model: {selected_model_name}")

    bundle = {
        "model": selected["model"],
        "encoder": encoder,
        "feature_names": feature_names,
        "selected_model": selected_model_name,
    }

    metadata = {
        "schemaVersion": config.export.schema_version,
        "datasetPath": str(dataset_path),
        "datasetRows": int(len(dataset)),
        "trainRows": int(len(train_frame)),
        "validationRows": int(len(validation_frame)),
        "selectedModel": selected_model_name,
        "selectionMetric": "log_loss",
        "trainedAt": datetime.now(timezone.utc).isoformat(),
        "trainDayStart": str(min(train_frame["event_day"])),
        "trainDayEnd": str(max(train_frame["event_day"])),
        "validationDayStart": str(min(validation_frame["event_day"])),
        "validationDayEnd": str(max(validation_frame["event_day"])),
        "validationDayPct": float(config.split.validation_day_pct),
        "validationDayCount": int(validation_frame["event_day"].nunique()),
        "dropDraws": bool(config.dataset.drop_draws),
    }

    metrics_payload = {
        "selectedModel": selected["metrics"],
        "allModels": {name: payload["metrics"] for name, payload in candidates.items()},
    }

    model_id, model_dir = export_model_bundle(
        config,
        run_context.models_dir,
        bundle,
        metadata,
        metrics_payload,
        feature_names,
    )
    update_registry(
        run_context.registry_file,
        model_id,
        model_dir,
        metadata,
        metrics_payload,
    )

    print(f"exported model: {model_dir}")
    print(f"updated registry: {run_context.registry_file}")
