from pathlib import Path

import pandas as pd

def write_parquet(frame: pd.DataFrame, path: Path) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)

    try:
        frame.to_parquet(path, index=False)
    except Exception as exc:
        raise RuntimeError(f"Failed to write file: {path}") from exc
