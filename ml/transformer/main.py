import argparse
from pathlib import Path

from config.config import load_config
from pipeline.build_aggregates import run_build_aggregates
from pipeline.build_dataset import run_build_dataset


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--config",
        default=str(Path(__file__).resolve().parent / "config" / "default.yml"),
        help="yaml config file",
    )
    return parser


def main() -> None:
    parser = build_parser()
    args = parser.parse_args()

    config = load_config(Path(args.config))

    run_build_aggregates(config)
    run_build_dataset(config)


if __name__ == "__main__":
    main()
