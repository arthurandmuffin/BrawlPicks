import math
import sys
from dataclasses import dataclass
from pathlib import Path

import joblib
import yaml

@dataclass
class LoadedModelBundle:
    model_id: str
    model_dir: Path
    model: object
    encoder: object
    feature_names: list[str]
    metadata: dict
    metrics: dict
    aggregate_dir: Path

class ModelBundleRepository:
    def __init__(self, fallback_aggregates_dir: Path, trainer_dir: Path):
        self.fallback_aggregates_dir = fallback_aggregates_dir
        self.trainer_dir = trainer_dir

    def load_bundle(self, model_dir: Path) -> LoadedModelBundle:
        #older joblib bundles still need trainer modules importable
        self._prepare_trainer_module_imports()

        bundle = joblib.load(model_dir / "model.joblib")
        metadata = self._read_yaml(model_dir / "metadata.yml")
        metrics = self._read_yaml(model_dir / "metrics.yml")
        feature_schema = self._read_yaml(model_dir / "feature_schema.yml")

        aggregate_dir = model_dir.parent.parent / "aggregates"
        if not aggregate_dir.exists():
            aggregate_dir = self.fallback_aggregates_dir

        return LoadedModelBundle(
            model_id=model_dir.name,
            model_dir=model_dir,
            model=bundle["model"],
            encoder=bundle["encoder"],
            feature_names=feature_schema.get("featureNames", bundle.get("feature_names", [])),
            metadata=metadata,
            metrics=metrics,
            aggregate_dir=aggregate_dir,
        )

    def predict_probability(self, bundle: LoadedModelBundle, feature_frame) -> float:
        model = bundle.model

        if hasattr(model, "predict_proba"):
            proba = model.predict_proba(feature_frame)
            return float(proba[0][1])

        if hasattr(model, "decision_function"):
            #fallback for models that expose scores instead of calibrated probs
            score = float(model.decision_function(feature_frame)[0])
            return 1.0 / (1.0 + math.exp(-score))

        raise ValueError("loaded model does not support predict_proba or decision_function")

    def _prepare_trainer_module_imports(self) -> None:
        trainer_dir = str(self.trainer_dir.resolve())
        if trainer_dir not in sys.path:
            sys.path.insert(0, trainer_dir)

    def _read_yaml(self, path: Path) -> dict:
        with path.open("r", encoding="utf-8") as handle:
            return yaml.safe_load(handle) or {}
