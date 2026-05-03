try:
    from transformer.config.config import Config
    from transformer.run_context import TransformerRunContext
    from transformer.pipeline.aggregates.builder import AggregateArtifacts
    from transformer.pipeline.battle.enricher import OUTPUT_COLUMNS, enrich_battles
    from transformer.pipeline.battle.reader import read_battles
    from transformer.pipeline.util.aggregate_paths import windowed_output_path
    from transformer.pipeline.util.date_windows import prior_window_for_day, target_dates
    from transformer.pipeline.util.parquet_reader import read_parquet
    from transformer.pipeline.util.parquet_writer import write_parquet
except ImportError:
    from config.config import Config
    from run_context import TransformerRunContext
    from pipeline.aggregates.builder import AggregateArtifacts
    from pipeline.battle.enricher import OUTPUT_COLUMNS, enrich_battles
    from pipeline.battle.reader import read_battles
    from pipeline.util.aggregate_paths import windowed_output_path
    from pipeline.util.date_windows import prior_window_for_day, target_dates
    from pipeline.util.parquet_reader import read_parquet
    from pipeline.util.parquet_writer import write_parquet

def run_build_dataset(config: Config, run_context: TransformerRunContext) -> None:
    target_days = target_dates(
        run_context.lookback_days,
        include_today=run_context.include_today,
    )
    enriched_frames = []

    for target_day in target_days:
        aggregate_start, aggregate_end = prior_window_for_day(
            target_day,
            config.aggregate_window.prior_days,
        )

        training_battles = read_battles(
            run_context.battle_logs_dir,
            target_day,
            target_day,
        )
        if training_battles.empty:
            continue

        aggregates = _load_aggregate_artifacts(
            config,
            run_context,
            aggregate_start.isoformat(),
            aggregate_end.isoformat(),
        )
        enriched = enrich_battles(
            training_battles,
            aggregates,
            config.features,
        )
        enriched_frames.append(enriched)

    if enriched_frames:
        import pandas as pd

        #concat stacks each per-day dataframe into one final dataset dataframe
        dataset = pd.concat(enriched_frames, ignore_index=True)
    else:
        import pandas as pd

        dataset = pd.DataFrame(columns=OUTPUT_COLUMNS)

    dataset_name = (
        f"training_rows_{target_days[0].isoformat()}_"
        f"{target_days[-1].isoformat()}.parquet"
    )
    write_parquet(
        dataset,
        run_context.datasets_dir / dataset_name,
    )


def _load_aggregate_artifacts(
    config: Config,
    run_context: TransformerRunContext,
    start_date: str,
    end_date: str,
) -> AggregateArtifacts:
    strength = read_parquet(
        windowed_output_path(
            run_context.aggregates_dir,
            config.paths.strength_file,
            start_date,
            end_date,
        )
    )
    synergy = read_parquet(
        windowed_output_path(
            run_context.aggregates_dir,
            config.paths.synergy_file,
            start_date,
            end_date,
        )
    )
    counter = read_parquet(
        windowed_output_path(
            run_context.aggregates_dir,
            config.paths.counter_file,
            start_date,
            end_date,
        )
    )

    return AggregateArtifacts(
        strength=strength,
        synergy=synergy,
        counter=counter,
    )
