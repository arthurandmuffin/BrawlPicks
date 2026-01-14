from dataclasses import dataclass
from pathlib import Path

@dataclass
class TransformerRunContext:
    battle_logs_dir: Path
    aggregates_dir: Path
    datasets_dir: Path
    lookback_days: int
    include_today: bool
