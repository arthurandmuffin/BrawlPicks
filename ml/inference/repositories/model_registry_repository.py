from pathlib import Path

import yaml

class ModelRegistryRepository:
    def __init__(self, models_dir: Path, registry_file: Path):
        self.models_dir = models_dir
        self.registry_file = registry_file

    def resolve_model_dir(self, model_id: str | None, model_dir: Path | None) -> Path:
        if model_dir is not None:
            return model_dir.resolve()

        if self.registry_file.exists():
            #prefer explicit registry choices before falling back to folder sorting
            registry = self._read_registry()
            if model_id is not None:
                for entry in registry.get("models", []):
                    if entry.get("modelId") == model_id:
                        return Path(entry["path"]).resolve()
                raise FileNotFoundError(f"model id not found in registry: {model_id}")

            latest_id = registry.get("latestCandidateModelId") or registry.get("activeModelId")
            if latest_id is not None:
                for entry in registry.get("models", []):
                    if entry.get("modelId") == latest_id:
                        return Path(entry["path"]).resolve()

        model_dirs = sorted(path for path in self.models_dir.iterdir() if path.is_dir())
        if not model_dirs:
            raise FileNotFoundError(f"no model directories found under {self.models_dir}")

        return model_dirs[-1].resolve()

    def _read_registry(self) -> dict:
        with self.registry_file.open("r", encoding="utf-8") as handle:
            return yaml.safe_load(handle) or {}
