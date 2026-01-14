from datetime import datetime, timezone
from pathlib import Path

import joblib

try:
    from trainer.config.config import Config
    from trainer.util.yaml_io import write_yaml_file
except ImportError:
    from config.config import Config
    from util.yaml_io import write_yaml_file

def export_model_bundle(
    config: Config,
    models_dir: Path,
    bundle: dict,
    metadata: dict,
    metrics: dict,
    feature_names: list[str],
) -> tuple[str, Path]:
    model_id = datetime.now(timezone.utc).strftime("model_%Y%m%dT%H%M%SZ")
    model_dir = models_dir / model_id
    model_dir.mkdir(parents=True, exist_ok=False)

    artifact_path = model_dir / config.export.artifact_file
    metadata_path = model_dir / config.export.metadata_file
    metrics_path = model_dir / config.export.metrics_file
    feature_schema_path = model_dir / config.export.feature_schema_file

    joblib.dump(bundle, artifact_path)

    write_yaml_file(metadata_path, metadata)
    write_yaml_file(metrics_path, metrics)
    write_yaml_file(
        feature_schema_path,
        {
            "schemaVersion": config.export.schema_version,
            "featureNames": feature_names,
        },
    )

    return model_id, model_dir
