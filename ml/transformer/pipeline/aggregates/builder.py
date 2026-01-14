from dataclasses import dataclass
from itertools import combinations
import math

import pandas as pd

try:
    from transformer.config.config import FeatureConfig
except ImportError:
    from config.config import FeatureConfig


@dataclass
class AggregateArtifacts:
    strength: pd.DataFrame
    synergy: pd.DataFrame
    counter: pd.DataFrame


def _normalize_team(value):
    if value is None:
        return []
    if isinstance(value, list):
        return [int(v) for v in value]
    #normalize parquet values into python list
    if hasattr(value, "tolist"):
        return [int(v) for v in value.tolist()]
    return [int(v) for v in value]


def _bucket_rank(rank_value, bucket_size):
    rank = int(rank_value)
    return (rank // bucket_size) * bucket_size


def _mean(values):
    return sum(values) / len(values)


def _variance(values, mean):
    return sum(math.pow(value - mean, 2) for value in values) / len(values)


def _avg_sampling_noise(probabilities, sample_sizes):
    noise_sum = 0.0
    for probability, sample_size in zip(probabilities, sample_sizes):
        noise_sum += probability * (1.0 - probability) / sample_size
    return noise_sum / len(probabilities)

#for baye shrink, similar to 3rd party upstream
def _variance_matched_k(probabilities, sample_sizes):
    mean = _mean(probabilities)
    avg_sampling_noise = _avg_sampling_noise(probabilities, sample_sizes)
    total_variance = _variance(probabilities, mean)
    skill_spread = total_variance - avg_sampling_noise

    if skill_spread <= 0:
        return 1e9, mean

    k = mean * (1.0 - mean) / skill_spread - 1.0
    if k < 0:
        return 0.0, mean

    return k, mean


def _apply_bayesian_shrink(grouped: pd.DataFrame) -> pd.DataFrame:
    if grouped.empty:
        grouped["score"] = []
        return grouped

    #tolist() to transform pandas math to use the mean / variance matching helpers above
    probabilities = (grouped["wins"] / grouped["total"]).tolist()
    sample_sizes = grouped["total"].tolist()
    k, mean = _variance_matched_k(probabilities, sample_sizes)

    #write entire column
    grouped["score"] = (
        (grouped["wins"] + (mean * k)) /
        (grouped["total"] + k)
    )
    return grouped


def build_aggregate_artifacts(
    battles: pd.DataFrame,
    features: FeatureConfig,
) -> AggregateArtifacts:
    strength_rows = []
    synergy_rows = []
    counter_rows = []

    for battle in battles.itertuples(index=False):
        rank_bucket = _bucket_rank(battle.rank, features.rank_bucket_size)
        team_w = _normalize_team(battle.team_W)
        team_l = _normalize_team(battle.team_L)
        draw_flag = bool(battle.draw_flag)

        teams = [
            (team_w, 1.0 if not draw_flag else 0.5),
            (team_l, 0.0 if not draw_flag else 0.5),
        ]

        for team, win_value in teams:
            for brawler_id in team:
                strength_rows.append(
                    {
                        "map_name": battle.map_name,
                        "mode": battle.mode,
                        "rank_bucket": rank_bucket,
                        "brawler_id": int(brawler_id),
                        "wins": win_value,
                        "total": 1.0,
                    }
                )

            # keep pair direction so no reliance on ID later, brittle possibly
            for left_id, right_id in combinations(team, 2):
                synergy_rows.append(
                    {
                        "map_name": battle.map_name,
                        "mode": battle.mode,
                        "rank_bucket": rank_bucket,
                        "left_id": int(left_id),
                        "right_id": int(right_id),
                        "wins": win_value,
                        "total": 1.0,
                    }
                )
                synergy_rows.append(
                    {
                        "map_name": battle.map_name,
                        "mode": battle.mode,
                        "rank_bucket": rank_bucket,
                        "left_id": int(right_id),
                        "right_id": int(left_id),
                        "wins": win_value,
                        "total": 1.0,
                    }
                )

        for winner_id in team_w:
            for loser_id in team_l:
                counter_rows.append(
                    {
                        "map_name": battle.map_name,
                        "mode": battle.mode,
                        "rank_bucket": rank_bucket,
                        "left_id": int(winner_id),
                        "right_id": int(loser_id),
                        "wins": 1.0 if not draw_flag else 0.5,
                        "total": 1.0,
                    }
                )
                counter_rows.append(
                    {
                        "map_name": battle.map_name,
                        "mode": battle.mode,
                        "rank_bucket": rank_bucket,
                        "left_id": int(loser_id),
                        "right_id": int(winner_id),
                        "wins": 0.0 if not draw_flag else 0.5,
                        "total": 1.0,
                    }
                )

    strength = _finalize_strength(strength_rows)
    synergy = _finalize_pairs(synergy_rows)
    counter = _finalize_pairs(counter_rows)

    return AggregateArtifacts(strength=strength, synergy=synergy, counter=counter)


def _finalize_strength(rows) -> pd.DataFrame:
    #dict rows into a real table for pandas
    frame = pd.DataFrame(rows)
    if frame.empty:
        return pd.DataFrame(
            columns=["map_name", "mode", "rank_bucket", "brawler_id", "wins", "total", "score"]
        )

    #basically GROUP BY keys, SUM(wins), SUM(total)
    grouped = (
        frame.groupby(["map_name", "mode", "rank_bucket", "brawler_id"], as_index=False)[["wins", "total"]].sum()
    )
    return _apply_bayesian_shrink(grouped)


def _finalize_pairs(rows) -> pd.DataFrame:
    frame = pd.DataFrame(rows)
    if frame.empty:
        return pd.DataFrame(
            columns=["map_name", "mode", "rank_bucket", "left_id", "right_id", "wins", "total", "score"]
        )

    #same as finalize_strength
    grouped = (
        frame.groupby(["map_name", "mode", "rank_bucket", "left_id", "right_id"], as_index=False)[["wins", "total"]].sum()
    )
    return _apply_bayesian_shrink(grouped)
