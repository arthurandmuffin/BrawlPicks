import argparse
from pathlib import Path

from config.config import load_config
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
    run_build_aggregates(config)
    run_build_dataset(config)


def main() -> None:
    parser = build_parser()
    args = parser.parse_args()
    run_from_config(Path(args.config))


if __name__ == "__main__":
    main()
