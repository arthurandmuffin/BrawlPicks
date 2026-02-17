from flask import Blueprint

from controllers.inference_controller import InferenceController

def build_blueprint(controller: InferenceController) -> Blueprint:
    blueprint = Blueprint("inference", __name__)

    #equivalent of routes.go
    blueprint.add_url_rule("/health", view_func=controller.health, methods=["GET"])
    blueprint.add_url_rule("/model", view_func=controller.model, methods=["GET"])
    blueprint.add_url_rule("/predict", view_func=controller.predict, methods=["POST"])
    blueprint.add_url_rule("/recommend", view_func=controller.recommend, methods=["POST"])

    return blueprint
