from datetime import date
from pathlib import Path
from typing import Iterable

import pandas as pd


REQUIRED_COLUMNS = [
    "timestamp",
    "map_name",
    "mode",
    "rank",
    "team_W",
    "team_L",
    "draw_flag",
]


def _date_range(start_date: date, end_date: date) -> Iterable[date]:
    current = start_date
    while current <= end_date:
        yield current
        current = current.fromordinal(current.toordinal() + 1)


def read_battles(root: Path, start_date: date, end_date: date) -> pd.DataFrame:
    frames = []
    for current in _date_range(start_date, end_date):
        day_dir = root / current.isoformat()
        if not day_dir.exists():
            continue

        for parquet_path in sorted(day_dir.rglob("*.parquet")):
            frame = pd.read_parquet(parquet_path)
            frames.append(frame)

    if not frames:
        return pd.DataFrame(columns=REQUIRED_COLUMNS)

    battles = pd.concat(frames, ignore_index=True)
    battles["timestamp"] = pd.to_datetime(battles["timestamp"], utc=True)
    return battles
