import argparse
from pathlib import Path

from app.app import create_app
from config.config import load_config

def default_config_path() -> Path:
    return Path(__file__).resolve().parent / "config" / "default.yml"

def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--config",
        default=str(default_config_path()),
        help="yaml config file",
    )
    parser.add_argument(
        "--model-id",
        default=None,
        help="specific model id to load from registry",
    )
    parser.add_argument(
        "--model-dir",
        default=None,
        help="specific model directory to load directly",
    )
    return parser

def main() -> None:
    parser = build_parser()
    args = parser.parse_args()

    config = load_config(Path(args.config))
    #pin one model by hand with model_dir
    model_dir = Path(args.model_dir).resolve() if args.model_dir else None
    app = create_app(config, model_id=args.model_id, model_dir=model_dir)
    app.run(host=config.server.host, port=config.server.port, debug=config.server.debug)

if __name__ == "__main__":
    main()
