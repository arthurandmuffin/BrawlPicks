from dataclasses import dataclass
from pathlib import Path

@dataclass
class TrainerRunContext:
    datasets_dir: Path
    models_dir: Path
    registry_file: Path
