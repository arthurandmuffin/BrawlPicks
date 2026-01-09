from sklearn.ensemble import HistGradientBoostingClassifier

from config.config import Config

def build_model(config: Config):
    return HistGradientBoostingClassifier(
        learning_rate=config.models.hist_gradient_boosting.learning_rate,
        max_depth=config.models.hist_gradient_boosting.max_depth,
        max_iter=config.models.hist_gradient_boosting.max_iter,
        random_state=17,
    )
