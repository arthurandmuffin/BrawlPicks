import pandas as pd

from repositories.aggregate_repository import AggregateRepository
from repositories.model_bundle_repository import LoadedModelBundle, ModelBundleRepository

class InferenceService:
    def __init__(
        self,
        model_bundle: LoadedModelBundle,
        model_repository: ModelBundleRepository,
        aggregate_repository: AggregateRepository,
        default_team_size: int,
        default_top_k: int,
    ):
        self.model_bundle = model_bundle
        self.model_repository = model_repository
        self.aggregate_repository = aggregate_repository
        self.default_team_size = default_team_size
        self.default_top_k = default_top_k
        self.known_brawler_ids = {int(value) for value in self.model_bundle.encoder.brawler_ids}

    def health(self) -> dict:
        return {
            "status": "ok",
            "model_id": self.model_bundle.model_id,
        }

    def model_info(self) -> dict:
        return {
            "model_id": self.model_bundle.model_id,
            "model_dir": str(self.model_bundle.model_dir),
            "metadata": self.model_bundle.metadata,
            "metrics": self.model_bundle.metrics,
            "aggregate_snapshots": self.aggregate_repository.snapshot_names,
        }

    def predict(self, map_name: str, mode: str, rank: int, team_a: list[int], team_b: list[int]) -> float:
        self._assert_known_brawlers("team_a", team_a)
        self._assert_known_brawlers("team_b", team_b)
        feature_row = self.aggregate_repository.build_feature_row(map_name, mode, rank, team_a, team_b)
        feature_frame = self.model_bundle.encoder.transform(pd.DataFrame([feature_row]))
        return self.model_repository.predict_probability(self.model_bundle, feature_frame)

    def recommend(
        self,
        map_name: str,
        mode: str,
        rank: int,
        ally_brawlers: list[int],
        enemy_brawlers: list[int],
        candidate_brawlers: list[int] | None,
        banned_brawlers: list[int],
        top_k: int | None,
    ) -> list[dict]:
        if len(ally_brawlers) >= self.default_team_size:
            raise ValueError("ally_brawlers already reached configured team size")
        self._assert_known_brawlers("ally_brawlers", ally_brawlers)
        self._assert_known_brawlers("enemy_brawlers", enemy_brawlers)
        self._assert_known_brawlers("banned_brawlers", banned_brawlers)
        if candidate_brawlers is not None:
            self._assert_known_brawlers("candidate_brawlers", candidate_brawlers)

        #score each legal next pick one by one by resulting wr
        candidates = self._candidate_pool(ally_brawlers, enemy_brawlers, candidate_brawlers, banned_brawlers)
        results = []

        for candidate in candidates:
            if len(ally_brawlers) < self.default_team_size:
                team_a = [*ally_brawlers, candidate]
            else:
                team_a = [*ally_brawlers]

            score = self.predict(map_name, mode, rank, team_a, enemy_brawlers)
            results.append(
                {
                    "brawler_id": int(candidate),
                    "score": score,
                }
            )

        results.sort(key=lambda item: item["score"], reverse=True)
        limit = top_k or self.default_top_k
        return results[:limit]

    def _candidate_pool(
        self,
        ally_brawlers: list[int],
        enemy_brawlers: list[int],
        candidate_brawlers: list[int] | None,
        banned_brawlers: list[int],
    ) -> list[int]:
        blocked = {int(value) for value in [*ally_brawlers, *enemy_brawlers, *banned_brawlers]}

        if candidate_brawlers is not None:
            pool = [int(value) for value in candidate_brawlers]
        else:
            #encoder vocab is our current known brawler universe
            pool = [int(value) for value in self.model_bundle.encoder.brawler_ids]

        return [candidate for candidate in pool if candidate not in blocked]

    def _assert_known_brawlers(self, field_name: str, brawler_ids: list[int]) -> None:
        unknown = sorted({int(value) for value in brawler_ids if int(value) not in self.known_brawler_ids})
        if unknown:
            raise ValueError(f"{field_name} contains unknown brawler ids: {unknown}")
