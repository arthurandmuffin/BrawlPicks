from flask import jsonify, request
from pydantic import ValidationError

from controllers.requests.inference_request import PredictRequest, RecommendRequest
from controllers.responses.inference_response import (
    health_response,
    model_response,
    predict_response,
    recommend_response,
)
from services.inference_service import InferenceService

class InferenceController:
    def __init__(self, service: InferenceService):
        self.service = service

    def health(self):
        return jsonify(health_response(self.service.health()))

    def model(self):
        return jsonify(model_response(self.service.model_info()))

    def predict(self):
        try:
            payload = PredictRequest.model_validate(request.get_json(force=True) or {})
        except ValidationError as exc:
            return jsonify({"error": "invalid request", "details": exc.errors()}), 400

        try:
            #service owns scoring logic
            score = self.service.predict(
                payload.map_name,
                payload.mode,
                payload.rank,
                payload.team_a,
                payload.team_b,
            )
        except ValueError as exc:
            return jsonify({"error": str(exc)}), 400

        return jsonify(predict_response(score, self.service.model_bundle.model_id))

    def recommend(self):
        try:
            payload = RecommendRequest.model_validate(request.get_json(force=True) or {})
        except ValidationError as exc:
            return jsonify({"error": "invalid request", "details": exc.errors()}), 400

        try:
            recommendations = self.service.recommend(
                payload.map_name,
                payload.mode,
                payload.rank,
                payload.ally_brawlers,
                payload.enemy_brawlers,
                payload.candidate_brawlers,
                payload.banned_brawlers,
                payload.top_k,
            )
        except ValueError as exc:
            return jsonify({"error": str(exc)}), 400

        return jsonify(recommend_response(self.service.model_bundle.model_id, recommendations))
