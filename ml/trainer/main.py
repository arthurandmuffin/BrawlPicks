import argparse
from pathlib import Path

try:
    from trainer.config.config import load_config
    from trainer.run_context import TrainerRunContext
    from trainer.pipeline.train import run_training
except ImportError:
    from config.config import load_config
    from run_context import TrainerRunContext
    from pipeline.train import run_training

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

def build_default_run_context(config) -> TrainerRunContext:
    return TrainerRunContext(
        datasets_dir=config.paths.datasets_dir,
        models_dir=config.paths.models_dir,
        registry_file=config.paths.registry_file,
    )

#orchestration enters here, skipping main
def run(config, run_context: TrainerRunContext) -> None:
    run_training(config, run_context)

def main() -> None:
    parser = build_parser()
    args = parser.parse_args()
    run_from_config(Path(args.config))

if __name__ == "__main__":
    main()
