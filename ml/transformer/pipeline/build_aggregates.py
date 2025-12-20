from config.config import Config
from pipeline.aggregates.builder import build_aggregate_artifacts
from pipeline.battle.reader import read_battles
from pipeline.util.aggregate_paths import windowed_output_path
from pipeline.util.date_windows import prior_window_for_day, target_dates
from pipeline.util.parquet_writer import write_parquet


def run_build_aggregates(config: Config) -> None:
    for target_day in target_dates(config.dataset.lookback_days):
        aggregate_start, aggregate_end = prior_window_for_day(
            target_day,
            config.aggregate_window.prior_days,
        )

        battles = read_battles(
            config.paths.battle_logs_dir,
            aggregate_start,
            aggregate_end,
        )
        aggregates = build_aggregate_artifacts(
            battles,
            config.features,
        )

        write_parquet(
            aggregates.strength,
            windowed_output_path(
                config.paths.aggregates_dir,
                config.paths.strength_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
        write_parquet(
            aggregates.synergy,
            windowed_output_path(
                config.paths.aggregates_dir,
                config.paths.synergy_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
        write_parquet(
            aggregates.counter,
            windowed_output_path(
                config.paths.aggregates_dir,
                config.paths.counter_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
