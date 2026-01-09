from dataclasses import dataclass

import pandas as pd

NUMERIC_COLUMNS = [
    "rank",
    "rank_bucket",
    "team_A_strength_avg",
    "team_B_strength_avg",
    "team_A_synergy_avg",
    "team_B_synergy_avg",
    "team_A_counter_avg",
    "team_B_counter_avg",
    "strength_delta",
    "synergy_delta",
    "counter_delta",
]

@dataclass
class FeatureEncoder:
    map_values: list[str]
    mode_values: list[str]
    brawler_ids: list[int]

    def feature_names(self) -> list[str]:
        names = list(NUMERIC_COLUMNS)
        names.extend(f"map::{value}" for value in self.map_values)
        names.extend(f"mode::{value}" for value in self.mode_values)
        names.extend(f"team_A::{brawler_id}" for brawler_id in self.brawler_ids)
        names.extend(f"team_B::{brawler_id}" for brawler_id in self.brawler_ids)
        return names

    def transform(self, frame: pd.DataFrame) -> pd.DataFrame:
        numeric = frame[NUMERIC_COLUMNS].reset_index(drop=True).astype(float)
        map_frame = _encode_categorical(frame["map_name"], "map", self.map_values).reset_index(drop=True)
        mode_frame = _encode_categorical(frame["mode"], "mode", self.mode_values).reset_index(drop=True)
        team_a_frame = _encode_team_lists(frame["team_A"], "team_A", self.brawler_ids).reset_index(drop=True)
        team_b_frame = _encode_team_lists(frame["team_B"], "team_B", self.brawler_ids).reset_index(drop=True)
        return pd.concat([numeric, map_frame, mode_frame, team_a_frame, team_b_frame], axis=1)


def fit_feature_encoder(frame: pd.DataFrame) -> FeatureEncoder:
    map_values = sorted(str(value) for value in frame["map_name"].dropna().unique())
    mode_values = sorted(str(value) for value in frame["mode"].dropna().unique())

    brawler_ids = set()
    for column in ["team_A", "team_B"]:
        for team in frame[column]:
            for brawler_id in _normalize_team(team):
                brawler_ids.add(int(brawler_id))

    return FeatureEncoder(map_values=map_values, mode_values=mode_values, brawler_ids=sorted(brawler_ids))


def _encode_categorical(series: pd.Series, prefix: str, allowed_values: list[str]) -> pd.DataFrame:
    data = {}
    as_strings = series.fillna("").astype(str)
    for value in allowed_values:
        data[f"{prefix}::{value}"] = (as_strings == value).astype(float)
    return pd.DataFrame(data)


def _encode_team_lists(series: pd.Series, prefix: str, brawler_ids: list[int]) -> pd.DataFrame:
    rows = []
    for team in series:
        team_ids = set(_normalize_team(team))
        rows.append({f"{prefix}::{brawler_id}": float(brawler_id in team_ids) for brawler_id in brawler_ids})
    return pd.DataFrame(rows)


def _normalize_team(value) -> list[int]:
    if value is None:
        return []
    if isinstance(value, list):
        return [int(v) for v in value]
    if hasattr(value, "tolist"):
        return [int(v) for v in value.tolist()]
    return [int(v) for v in value]
