from itertools import combinations
from pathlib import Path

import pandas as pd

class AggregateRepository:
    def __init__(self, aggregate_dir: Path, rank_bucket_size: int):
        self.aggregate_dir = aggregate_dir
        self.rank_bucket_size = rank_bucket_size
        self.strength_lookup = {}
        self.synergy_lookup = {}
        self.counter_lookup = {}
        self.snapshot_names = {}
        self._load_latest_snapshots()

    def build_feature_row(self, map_name: str, mode: str, rank: int, team_a: list[int], team_b: list[int]) -> dict:
        #this mirrors the trainer-side transformed row shape
        rank_bucket = self._bucket_rank(rank)
        context = (map_name, mode, rank_bucket)

        team_a_strength = self._average_strength(team_a, context)
        team_b_strength = self._average_strength(team_b, context)
        team_a_synergy = self._average_synergy(team_a, context)
        team_b_synergy = self._average_synergy(team_b, context)
        team_a_counter = self._average_counter(team_a, team_b, context)
        team_b_counter = self._average_counter(team_b, team_a, context)

        return {
            "map_name": map_name,
            "mode": mode,
            "rank": int(rank),
            "rank_bucket": rank_bucket,
            "team_A": [int(value) for value in team_a],
            "team_B": [int(value) for value in team_b],
            "team_A_strength_avg": team_a_strength,
            "team_B_strength_avg": team_b_strength,
            "team_A_synergy_avg": team_a_synergy,
            "team_B_synergy_avg": team_b_synergy,
            "team_A_counter_avg": team_a_counter,
            "team_B_counter_avg": team_b_counter,
            "strength_delta": team_a_strength - team_b_strength,
            "synergy_delta": team_a_synergy - team_b_synergy,
            "counter_delta": team_a_counter - team_b_counter,
        }

    def _load_latest_snapshots(self) -> None:
        strength_path = self._latest_snapshot("strength")
        synergy_path = self._latest_snapshot("synergy")
        counter_path = self._latest_snapshot("counter")

        self.snapshot_names = {
            "strength": strength_path.name,
            "synergy": synergy_path.name,
            "counter": counter_path.name,
        }

        #load the newest aggregate set found for this run/source
        strength = pd.read_parquet(strength_path)
        synergy = pd.read_parquet(synergy_path)
        counter = pd.read_parquet(counter_path)

        self.strength_lookup = {
            (row.map_name, row.mode, int(row.rank_bucket), int(row.brawler_id)): float(row.score)
            for row in strength.itertuples(index=False)
        }
        self.synergy_lookup = {
            (row.map_name, row.mode, int(row.rank_bucket), int(row.left_id), int(row.right_id)): float(row.score)
            for row in synergy.itertuples(index=False)
        }
        self.counter_lookup = {
            (row.map_name, row.mode, int(row.rank_bucket), int(row.left_id), int(row.right_id)): float(row.score)
            for row in counter.itertuples(index=False)
        }

    def _latest_snapshot(self, prefix: str) -> Path:
        matches = sorted(self.aggregate_dir.glob(f"{prefix}_*.parquet"))
        if not matches:
            raise FileNotFoundError(f"no aggregate snapshots found for {prefix} under {self.aggregate_dir}")
        return matches[-1]

    def _bucket_rank(self, rank: int) -> int:
        return (int(rank) // self.rank_bucket_size) * self.rank_bucket_size

    def _average_strength(self, team: list[int], context: tuple) -> float:
        if not team:
            return 0.0

        scores = [self.strength_lookup.get((context[0], context[1], context[2], brawler_id), 0.5) for brawler_id in team]
        return sum(scores) / len(scores)

    def _average_synergy(self, team: list[int], context: tuple) -> float:
        pairs = list(combinations(team, 2))
        if not pairs:
            return 0.5

        scores = [
            self.synergy_lookup.get((context[0], context[1], context[2], left_id, right_id), 0.5)
            for left_id, right_id in pairs
        ]
        return sum(scores) / len(scores)

    def _average_counter(self, source_team: list[int], target_team: list[int], context: tuple) -> float:
        if not source_team or not target_team:
            return 0.5

        #average every all source v. target pair into one matchup score
        scores = []
        for left_id in source_team:
            for right_id in target_team:
                scores.append(
                    self.counter_lookup.get((context[0], context[1], context[2], left_id, right_id), 0.5)
                )

        return sum(scores) / len(scores)
