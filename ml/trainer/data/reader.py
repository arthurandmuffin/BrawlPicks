from pathlib import Path

import pandas as pd

def load_latest_dataset(datasets_dir: Path, dataset_glob: str) -> tuple[pd.DataFrame, Path]:
    files = sorted(datasets_dir.glob(dataset_glob))
    if not files:
        raise FileNotFoundError(f"no dataset files found in {datasets_dir} for glob {dataset_glob}")

    # walk newest-first, but skip empty artifacts so a bad recent run does not poison trainer
    for dataset_path in reversed(files):
        frame = pd.read_parquet(dataset_path)
        if frame.empty:
            continue

        frame["timestamp"] = pd.to_datetime(frame["timestamp"], utc=True)
        frame["event_day"] = frame["timestamp"].dt.date
        return frame, dataset_path

    raise ValueError(f"all dataset files in {datasets_dir} are empty for glob {dataset_glob}")
