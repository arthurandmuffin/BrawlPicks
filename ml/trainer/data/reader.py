from pathlib import Path

import pandas as pd

def load_latest_dataset(datasets_dir: Path, dataset_glob: str) -> tuple[pd.DataFrame, Path]:
    files = sorted(datasets_dir.glob(dataset_glob))
    if not files:
        raise FileNotFoundError(f"no dataset files found in {datasets_dir} for glob {dataset_glob}")

    dataset_path = files[-1]
    frame = pd.read_parquet(dataset_path)
    frame["timestamp"] = pd.to_datetime(frame["timestamp"], utc=True)
    frame["event_day"] = frame["timestamp"].dt.date
    return frame, dataset_path
