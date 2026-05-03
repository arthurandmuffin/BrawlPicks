from pathlib import Path


def windowed_output_path(root: Path, filename: str, start_date: str, end_date: str) -> Path:
    base = Path(filename)
    stem = base.stem
    suffix = base.suffix or ".parquet"
    output_name = f"{stem}_{start_date}_{end_date}{suffix}"
    return root / output_name
