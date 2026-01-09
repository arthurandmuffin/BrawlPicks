import math

from sklearn.metrics import accuracy_score, brier_score_loss, log_loss, roc_auc_score

def evaluate_binary_classifier(y_true, y_proba) -> dict:
    #shift values away from exactly 1/0 to avoid infinity
    clipped = [min(max(float(value), 1e-9), 1.0 - 1e-9) for value in y_proba]
    labels = [1 if value >= 0.5 else 0 for value in clipped]

    metrics = {
        "accuracy": float(accuracy_score(y_true, labels)),
        "brier_score": float(brier_score_loss(y_true, clipped)),
        "log_loss": float(log_loss(y_true, clipped)),
    }

    if len(set(y_true)) > 1:
        metrics["roc_auc"] = float(roc_auc_score(y_true, clipped))
    else:
        metrics["roc_auc"] = math.nan

    return metrics
