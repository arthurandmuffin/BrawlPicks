from itertools import combinations

import pandas as pd

from config.config import FeatureConfig
from pipeline.aggregates.builder import AggregateArtifacts


def _normalize_team(value):
    if value is None:
        return []
    if isinstance(value, list):
        return [int(v) for v in value]
    #parquet array to py list
    if hasattr(value, "tolist"):
        return [int(v) for v in value.tolist()]
    return [int(v) for v in value]


def _bucket_rank(rank_value, bucket_size):
    rank = int(rank_value)
    return (rank // bucket_size) * bucket_size


def enrich_battles(
    battles: pd.DataFrame,
    aggregates: AggregateArtifacts,
    features: FeatureConfig,
) -> pd.DataFrame:
    #convert aggregate tables into dict lookups 
    strength_lookup = _build_strength_lookup(aggregates.strength)
    synergy_lookup = _build_pair_lookup(aggregates.synergy)
    counter_lookup = _build_pair_lookup(aggregates.counter)

    rows = []

    for battle in battles.itertuples(index=False):
        team_w = _normalize_team(battle.team_W)
        team_l = _normalize_team(battle.team_L)
        rank_bucket = _bucket_rank(battle.rank, features.rank_bucket_size)
        context = (battle.map_name, battle.mode, rank_bucket)
        draw_flag = bool(battle.draw_flag)

        winner_strength = _average_strength(team_w, context, strength_lookup)
        loser_strength = _average_strength(team_l, context, strength_lookup)
        winner_synergy = _average_synergy(team_w, context, synergy_lookup)
        loser_synergy = _average_synergy(team_l, context, synergy_lookup)
        winner_counter = _average_counter(team_w, team_l, context, counter_lookup)
        loser_counter = _average_counter(team_l, team_w, context, counter_lookup)

        base_row = {
            "timestamp": battle.timestamp,
            "map_name": battle.map_name,
            "mode": battle.mode,
            "rank": int(battle.rank),
            "rank_bucket": rank_bucket,
        }

        rows.append(
            {
                **base_row,
                "team_A": team_w,
                "team_B": team_l,
                "team_A_strength_avg": winner_strength,
                "team_B_strength_avg": loser_strength,
                "team_A_synergy_avg": winner_synergy,
                "team_B_synergy_avg": loser_synergy,
                "team_A_counter_avg": winner_counter,
                "team_B_counter_avg": loser_counter,
                "strength_delta": winner_strength - loser_strength,
                "synergy_delta": winner_synergy - loser_synergy,
                "counter_delta": winner_counter - loser_counter,
                "team_A_won": 0.5 if draw_flag else 1.0,
            }
        )

        #duplicate battle so (A, B, W) (B, A, L), for better learning (hope)
        rows.append(
            {
                **base_row,
                "team_A": team_l,
                "team_B": team_w,
                "team_A_strength_avg": loser_strength,
                "team_B_strength_avg": winner_strength,
                "team_A_synergy_avg": loser_synergy,
                "team_B_synergy_avg": winner_synergy,
                "team_A_counter_avg": loser_counter,
                "team_B_counter_avg": winner_counter,
                "strength_delta": loser_strength - winner_strength,
                "synergy_delta": loser_synergy - winner_synergy,
                "counter_delta": loser_counter - winner_counter,
                "team_A_won": 0.5 if draw_flag else 0.0,
            }
        )

    #[dict] -> dataframe conversion for parquet write
    return pd.DataFrame(rows)


def _build_strength_lookup(frame: pd.DataFrame):
    lookup = {}
    #itertuples walks dataframe rows without extra pandas index
    for row in frame.itertuples(index=False):
        lookup[(row.map_name, row.mode, int(row.rank_bucket), int(row.brawler_id))] = float(row.score)
    return lookup


def _build_pair_lookup(frame: pd.DataFrame):
    lookup = {}
    for row in frame.itertuples(index=False):
        lookup[(row.map_name, row.mode, int(row.rank_bucket), int(row.left_id), int(row.right_id))] = float(row.score)
    return lookup


def _average_strength(team, context, lookup):
    if not team:
        return 0.0

    #default to 0.5 when no aggregate history
    scores = [lookup.get((context[0], context[1], context[2], brawler_id), 0.5) for brawler_id in team]
    return sum(scores) / len(scores)


def _average_synergy(team, context, lookup):
    pairs = list(combinations(team, 2))
    if not pairs:
        return 0.5

    scores = [
        lookup.get((context[0], context[1], context[2], left_id, right_id), 0.5)
        for left_id, right_id in pairs
    ]
    return sum(scores) / len(scores)


def _average_counter(source_team, target_team, context, lookup):
    if not source_team or not target_team:
        return 0.5

    scores = []
    for left_id in source_team:
        for right_id in target_team:
            scores.append(lookup.get((context[0], context[1], context[2], left_id, right_id), 0.5))

    return sum(scores) / len(scores)
