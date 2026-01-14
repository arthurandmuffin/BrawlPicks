import argparse
from datetime import datetime
from pathlib import Path

from config.config import load_config
from trainer.config.config import load_config as load_trainer_config
from trainer.main import run as run_trainer
from trainer.run_context import TrainerRunContext
from transformer.config.config import load_config as load_transformer_config
from transformer.main import run as run_transformer
from transformer.run_context import TransformerRunContext

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
    run_id = datetime.now().strftime("%Y-%m-%d-%H%M%S")
    run_dir = config.paths.artifacts_dir / run_id

    print(f"run ID: {run_id}")
    print(f"run Dir: {run_dir}")

    transformer_config = load_transformer_config(config.paths.transformer_config)
    #override transformer internal configs
    transformer_run_context = TransformerRunContext(
        battle_logs_dir=config.paths.battle_logs_dir,
        aggregates_dir=run_dir / "aggregates",
        datasets_dir=run_dir / "datasets",
        lookback_days=config.run.lookback_days,
        include_today=config.run.include_today,
    )
    run_transformer(transformer_config, transformer_run_context)

    trainer_config = load_trainer_config(config.paths.trainer_config)
    #override trainer internal configs
    trainer_run_context = TrainerRunContext(
        datasets_dir=run_dir / "datasets",
        models_dir=run_dir / "models",
        registry_file=run_dir / "models" / "registry.yml",
    )
    run_trainer(trainer_config, trainer_run_context)

def main() -> None:
    parser = build_parser()
    args = parser.parse_args()
    run_from_config(Path(args.config))

if __name__ == "__main__":
    main()