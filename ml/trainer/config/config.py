from dataclasses import dataclass
from pathlib import Path

from util.pathing import resolve_path
from util.yaml_io import read_yaml_file


@dataclass
class PathConfig:
    output_dir: Path
    datasets_dir: Path
    models_dir: Path
    registry_file: Path


@dataclass
class DatasetConfig:
    dataset_glob: str
    drop_draws: bool


@dataclass
class SplitConfig:
    validation_day_pct: float


@dataclass
class LogisticRegressionConfig:
    c: float
    max_iter: int


@dataclass
class HistGradientBoostingConfig:
    learning_rate: float
    max_depth: int
    max_iter: int


@dataclass
class ModelConfig:
    logistic_regression: LogisticRegressionConfig
    hist_gradient_boosting: HistGradientBoostingConfig


@dataclass
class ExportConfig:
    schema_version: int
    artifact_file: str
    metadata_file: str
    metrics_file: str
    feature_schema_file: str


@dataclass
class Config:
    paths: PathConfig
    dataset: DatasetConfig
    split: SplitConfig
    models: ModelConfig
    export: ExportConfig


def load_config(path: Path) -> Config:
    raw = read_yaml_file(path)
    output_dir = Path(raw["paths"]["outputDir"]).resolve()

    config = Config(
        paths=PathConfig(
            output_dir=output_dir,
            datasets_dir=resolve_path(output_dir, raw["paths"]["datasetsDir"]),
            models_dir=resolve_path(output_dir, raw["paths"]["modelsDir"]),
            registry_file=resolve_path(output_dir, raw["paths"]["registryFile"]),
        ),
        dataset=DatasetConfig(
            dataset_glob=raw["dataset"]["datasetGlob"],
            drop_draws=bool(raw["dataset"]["dropDraws"]),
        ),
        split=SplitConfig(
            validation_day_pct=float(raw["split"]["validationDayPct"]),
        ),
        models=ModelConfig(
            logistic_regression=LogisticRegressionConfig(
                c=float(raw["models"]["logisticRegression"]["c"]),
                max_iter=int(raw["models"]["logisticRegression"]["maxIter"]),
            ),
            hist_gradient_boosting=HistGradientBoostingConfig(
                learning_rate=float(raw["models"]["histGradientBoosting"]["learningRate"]),
                max_depth=int(raw["models"]["histGradientBoosting"]["maxDepth"]),
                max_iter=int(raw["models"]["histGradientBoosting"]["maxIter"]),
            ),
        ),
        export=ExportConfig(
            schema_version=int(raw["export"]["schemaVersion"]),
            artifact_file=raw["export"]["artifactFile"],
            metadata_file=raw["export"]["metadataFile"],
            metrics_file=raw["export"]["metricsFile"],
            feature_schema_file=raw["export"]["featureSchemaFile"],
        ),
    )

    if not 0 < config.split.validation_day_pct < 1:
        raise ValueError("split.validationDayPct must be between 0 and 1")

    return config
