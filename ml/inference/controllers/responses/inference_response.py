def health_response(payload: dict) -> dict:
    return payload


def model_response(payload: dict) -> dict:
    return payload


def predict_response(score: float, model_id: str) -> dict:
    return {
        "model_id": model_id,
        "team_a_win_probability": score,
    }


def recommend_response(model_id: str, recommendations: list[dict]) -> dict:
    return {
        "model_id": model_id,
        "recommendations": recommendations,
    }
