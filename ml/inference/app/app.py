from flask import Flask

from config.config import Config
from controllers.inference_controller import InferenceController
from repositories.aggregate_repository import AggregateRepository
from repositories.model_bundle_repository import ModelBundleRepository
from repositories.model_registry_repository import ModelRegistryRepository
from routes.routes import build_blueprint
from services.inference_service import InferenceService

def create_app(config: Config, model_id: str | None = None, model_dir=None) -> Flask:
    # registry picks the bundle, then the rest of the app hangs off that
    registry_repository = ModelRegistryRepository(
        config.paths.models_dir,
        config.paths.registry_file,
    )
    resolved_model_dir = registry_repository.resolve_model_dir(model_id, model_dir)

    model_repository = ModelBundleRepository(
        config.paths.aggregates_dir,
        config.paths.trainer_dir,
    )
    model_bundle = model_repository.load_bundle(resolved_model_dir)

    aggregate_repository = AggregateRepository(
        model_bundle.aggregate_dir,
        config.features.rank_bucket_size,
    )

    service = InferenceService(
        model_bundle,
        model_repository,
        aggregate_repository,
        config.features.default_team_size,
        config.features.default_top_k,
    )

    controller = InferenceController(service)

    app = Flask(__name__)
    app.register_blueprint(build_blueprint(controller))
    return app
