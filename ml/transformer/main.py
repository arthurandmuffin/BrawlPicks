import argparse
from pathlib import Path

#different import depending on direct run vs orchestration from ml/main.py
try:
    from transformer.config.config import load_config
    from transformer.run_context import TransformerRunContext
    from transformer.pipeline.build_aggregates import run_build_aggregates
    from transformer.pipeline.build_dataset import run_build_dataset
except ImportError:
    from config.config import load_config
    from run_context import TransformerRunContext
    from pipeline.build_aggregates import run_build_aggregates
    from pipeline.build_dataset import run_build_dataset

def default_config_path() -> Path:
    return Path(__file__).resolve().parent / "config" / "default.yml"

def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--config",
        default=str(default_config_path()),
        help="yaml config file",
    )
    return parser

def run_from_config(config_path: Path) -> None:
    config = load_config(config_path)
    run(config, build_default_run_context(config))

def build_default_run_context(config) -> TransformerRunContext:
    return TransformerRunContext(
        battle_logs_dir=config.paths.battle_logs_dir,
        aggregates_dir=config.paths.aggregates_dir,
        datasets_dir=config.paths.datasets_dir,
        lookback_days=config.dataset.lookback_days,
        include_today=config.dataset.include_today,
    )

#orchestration enters here, skipping main
def run(config, run_context: TransformerRunContext) -> None:
    run_build_aggregates(config, run_context)
    run_build_dataset(config, run_context)

def main() -> None:
    parser = build_parser()
    args = parser.parse_args()
    run_from_config(Path(args.config))

if __name__ == "__main__":
    main()
