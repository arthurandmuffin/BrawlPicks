from pathlib import Path

import pandas as pd

def read_parquet(path: Path) -> pd.DataFrame:
    if not path.exists():
        raise FileNotFoundError(f"parquet file not found: {path}")
    return pd.read_parquet(path)
