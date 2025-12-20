from dataclasses import dataclass
from pathlib import Path

import yaml


@dataclass
class PathConfig:
    source_dir: Path
    output_dir: Path
    battle_logs_dir: Path
    scraper_synergy_dir: Path
    aggregates_dir: Path
    datasets_dir: Path
    strength_file: str
    synergy_file: str
    counter_file: str


@dataclass
class DatasetConfig:
    lookback_days: int


@dataclass
class AggregateWindowConfig:
    prior_days: int


@dataclass
class FeatureConfig:
    rank_bucket_size: int
    include_draws: bool


@dataclass
class Config:
    paths: PathConfig
    dataset: DatasetConfig
    aggregate_window: AggregateWindowConfig
    features: FeatureConfig


def load_config(path: Path) -> Config:
    with path.open("r", encoding="utf-8") as handle:
        raw = yaml.safe_load(handle)

    source_dir = Path(raw["paths"]["sourceDir"]).resolve()
    output_dir = Path(raw["paths"]["outputDir"]).resolve()

    config = Config(
        paths=PathConfig(
            source_dir=source_dir,
            output_dir=output_dir,
            battle_logs_dir=_resolve_path(source_dir, raw["paths"]["battleLogsDir"]),
            scraper_synergy_dir=_resolve_path(source_dir, raw["paths"]["scraperSynergyDir"]),
            aggregates_dir=_resolve_path(output_dir, raw["paths"]["aggregatesDir"]),
            datasets_dir=_resolve_path(output_dir, raw["paths"]["datasetsDir"]),
            strength_file=raw["paths"]["strengthFile"],
            synergy_file=raw["paths"]["synergyFile"],
            counter_file=raw["paths"]["counterFile"],
        ),
        dataset=DatasetConfig(
            lookback_days=int(raw["dataset"]["lookbackDays"]),
        ),
        aggregate_window=AggregateWindowConfig(
            prior_days=int(raw["aggregateWindow"]["priorDays"]),
        ),
        features=FeatureConfig(
            rank_bucket_size=int(raw["features"]["rankBucketSize"]),
            include_draws=bool(raw["features"]["includeDraws"]),
        ),
    )

    if config.dataset.lookback_days <= 0:
        print("dataset.lookbakcDays set to 1")
        config.dataset.lookback_days = 1

    if config.aggregate_window.prior_days <= 0:
        print("aggregateWindow.priorDays set to 1")
        config.aggregate_window.prior_days = 1

    if config.features.rank_bucket_size <= 0:
        print("rankBucketSize defaulting to 1")
        config.features.rank_bucket_size = 1

    return config


def _resolve_path(base_dir: Path, value: str) -> Path:
    path = Path(value)
    if path.is_absolute():
        return path
    return (base_dir / path).resolve()
