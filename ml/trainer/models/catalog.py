try:
    from trainer.models.sklearn import hist_gradient_boosting, logistic_regression
except ImportError:
    from models.sklearn import hist_gradient_boosting, logistic_regression

def build_candidate_models(config) -> list[tuple[str, object]]:
    return [
        ("logistic_regression", logistic_regression.build_model(config)),
        ("hist_gradient_boosting", hist_gradient_boosting.build_model(config)),
    ]
