try:
    from transformer.config.config import Config
    from transformer.run_context import TransformerRunContext
    from transformer.pipeline.aggregates.builder import build_aggregate_artifacts
    from transformer.pipeline.battle.reader import read_battles
    from transformer.pipeline.util.aggregate_paths import windowed_output_path
    from transformer.pipeline.util.date_windows import prior_window_for_day, target_dates
    from transformer.pipeline.util.parquet_writer import write_parquet
except ImportError:
    from config.config import Config
    from run_context import TransformerRunContext
    from pipeline.aggregates.builder import build_aggregate_artifacts
    from pipeline.battle.reader import read_battles
    from pipeline.util.aggregate_paths import windowed_output_path
    from pipeline.util.date_windows import prior_window_for_day, target_dates
    from pipeline.util.parquet_writer import write_parquet

def run_build_aggregates(config: Config, run_context: TransformerRunContext) -> None:
    for target_day in target_dates(
        run_context.lookback_days,
        include_today=run_context.include_today,
    ):
        aggregate_start, aggregate_end = prior_window_for_day(
            target_day,
            config.aggregate_window.prior_days,
        )

        battles = read_battles(
            run_context.battle_logs_dir,
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
                run_context.aggregates_dir,
                config.paths.strength_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
        write_parquet(
            aggregates.synergy,
            windowed_output_path(
                run_context.aggregates_dir,
                config.paths.synergy_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
        write_parquet(
            aggregates.counter,
            windowed_output_path(
                run_context.aggregates_dir,
                config.paths.counter_file,
                aggregate_start.isoformat(),
                aggregate_end.isoformat(),
            ),
        )
