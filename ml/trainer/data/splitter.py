import math

import pandas as pd

def split_train_validation(frame: pd.DataFrame, validation_day_pct: float) -> tuple[pd.DataFrame, pd.DataFrame]:
    unique_days = sorted(frame["event_day"].unique())
    if len(unique_days) < 2:
        raise ValueError("need at least 2 unique days to create a train/validation split")

    validation_days = round(len(unique_days) * validation_day_pct)
    validation_days = max(1, validation_days)
    validation_days = min(validation_days, len(unique_days) - 1)

    validation_day_set = set(unique_days[-validation_days:])
    #rows where day not in validation set
    train = frame[~frame["event_day"].isin(validation_day_set)].copy()
    #rows in validation set
    validation = frame[frame["event_day"].isin(validation_day_set)].copy()

    if train.empty or validation.empty:
        raise ValueError("train/validation split failed")

    return train, validation
