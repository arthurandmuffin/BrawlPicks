from pathlib import Path

from config.config import load_config
from pipeline.train import run_training

def main() -> None:
    config_path = Path(__file__).resolve().parent / "config" / "default.yml"
    config = load_config(config_path)
    run_training(config)


if __name__ == "__main__":
    main()
