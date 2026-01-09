from datetime import datetime, timezone
from pathlib import Path

from util.yaml_io import read_yaml_file, write_yaml_file

def update_registry(
    registry_path: Path,
    model_id: str,
    model_dir: Path,
    metadata: dict,
    metrics: dict,
) -> None:
    registry_path.parent.mkdir(parents=True, exist_ok=True)

    if registry_path.exists():
        registry = read_yaml_file(registry_path)
    else:
        registry = {}

    entries = registry.get("models", [])
    entries.append(
        {
            "modelId": model_id,
            "path": str(model_dir),
            "selectedModel": metadata["selectedModel"],
            "datasetPath": metadata["datasetPath"],
            "createdAt": datetime.now(timezone.utc).isoformat(),
            "validationMetrics": metrics["selectedModel"],
            "status": "candidate",
        }
    )

    registry["models"] = entries
    registry["latestCandidateModelId"] = model_id
    if "activeModelId" not in registry:
        registry["activeModelId"] = model_id

    write_yaml_file(registry_path, registry)
