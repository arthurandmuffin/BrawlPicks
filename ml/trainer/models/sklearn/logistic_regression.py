from sklearn.linear_model import LogisticRegression

try:
    from trainer.config.config import Config
except ImportError:
    from config.config import Config

def build_model(config: Config):
    return LogisticRegression(
        C=config.models.logistic_regression.c,
        max_iter=config.models.logistic_regression.max_iter,
        random_state=17,
    )
