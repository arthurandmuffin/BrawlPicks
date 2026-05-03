from dataclasses import dataclass
from pathlib import Path

import yaml

@dataclass
class PathConfig:
    artifacts_dir: Path
    battle_logs_dir: Path
    transformer_config: Path
    trainer_config: Path

@dataclass
class RunConfig:
    lookback_days: int
    include_today: bool

@dataclass
class Config:
    paths: PathConfig
    run: RunConfig


def load_config(path: Path) -> Config:
    with path.open("r", encoding="utf-8") as handle:
        raw = yaml.safe_load(handle) or {}
    base_dir = path.resolve().parent.parent

    config = Config(
        paths=PathConfig(
            artifacts_dir=_resolve_path(base_dir, raw["paths"]["artifactsDir"]),
            battle_logs_dir=_resolve_path(base_dir, raw["paths"]["battleLogsDir"]),
            transformer_config=_resolve_path(base_dir, raw["paths"]["transformerConfig"]),
            trainer_config=_resolve_path(base_dir, raw["paths"]["trainerConfig"]),
        ),
        run=RunConfig(
            lookback_days=int(raw["run"]["lookbackDays"]),
            include_today=bool(raw["run"]["includeToday"]),
        ),
    )

    if config.run.lookback_days <= 0:
        print("run.lookbackDays defaulted to 1")
        config.run.lookback_days = 1

    return config


def _resolve_path(base_dir: Path, value: str) -> Path:
    path = Path(value)
    if path.is_absolute():
        return path
    return (base_dir / path).resolve()
