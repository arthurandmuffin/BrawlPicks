from dataclasses import dataclass
from pathlib import Path

import yaml

@dataclass
class PathConfig:
    output_dir: Path
    models_dir: Path
    registry_file: Path
    aggregates_dir: Path
    trainer_dir: Path

@dataclass
class ServerConfig:
    host: str
    port: int
    debug: bool

@dataclass
class FeatureConfig:
    rank_bucket_size: int
    default_team_size: int
    default_top_k: int

@dataclass
class Config:
    paths: PathConfig
    server: ServerConfig
    features: FeatureConfig

def load_config(path: Path) -> Config:
    with path.open("r", encoding="utf-8") as handle:
        raw = yaml.safe_load(handle) or {}

    output_dir = Path(raw["paths"]["outputDir"]).resolve()

    config = Config(
        paths=PathConfig(
            output_dir=output_dir,
            models_dir=_resolve_path(output_dir, raw["paths"]["modelsDir"]),
            registry_file=_resolve_path(output_dir, raw["paths"]["registryFile"]),
            aggregates_dir=_resolve_path(output_dir, raw["paths"]["aggregatesDir"]),
            trainer_dir=_resolve_path(output_dir, raw["paths"]["trainerDir"]),
        ),
        server=ServerConfig(
            host=raw["server"]["host"],
            port=int(raw["server"]["port"]),
            debug=bool(raw["server"]["debug"]),
        ),
        features=FeatureConfig(
            rank_bucket_size=int(raw["features"]["rankBucketSize"]),
            default_team_size=int(raw["features"]["defaultTeamSize"]),
            default_top_k=int(raw["features"]["defaultTopK"]),
        ),
    )

    if config.server.port <= 0:
        raise ValueError("server.port must be positive")

    if config.features.rank_bucket_size <= 0:
        raise ValueError("features.rankBucketSize must be positive")

    if config.features.default_team_size <= 0:
        raise ValueError("features.defaultTeamSize must be positive")

    if config.features.default_top_k <= 0:
        raise ValueError("features.defaultTopK must be positive")

    return config

def _resolve_path(base_dir: Path, value: str) -> Path:
    path = Path(value)
    if path.is_absolute():
        return path
    return (base_dir / path).resolve()
