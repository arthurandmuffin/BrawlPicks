from config.config import Config
from pipeline.aggregates.builder import AggregateArtifacts
from pipeline.battle.enricher import enrich_battles
from pipeline.battle.reader import read_battles
from pipeline.util.aggregate_paths import windowed_output_path
from pipeline.util.date_windows import prior_window_for_day, target_dates
from pipeline.util.parquet_reader import read_parquet
from pipeline.util.parquet_writer import write_parquet


def run_build_dataset(config: Config) -> None:
    target_days = target_dates(config.dataset.lookback_days)
    enriched_frames = []

    for target_day in target_days:
        aggregate_start, aggregate_end = prior_window_for_day(
            target_day,
            config.aggregate_window.prior_days,
        )

        training_battles = read_battles(
            config.paths.battle_logs_dir,
            target_day,
            target_day,
        )
        if training_battles.empty:
            continue

        aggregates = _load_aggregate_artifacts(
            config,
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
        training_battles = read_battles(
            config.paths.battle_logs_dir,
            target_days[0],
            target_days[-1],
        )
        #keep the empty result as a dataframe with the same base shape
        dataset = training_battles.iloc[0:0].copy()

    dataset_name = (
        f"training_rows_{target_days[0].isoformat()}_"
        f"{target_days[-1].isoformat()}.parquet"
    )
    write_parquet(
        dataset,
        config.paths.datasets_dir / dataset_name,
    )


def _load_aggregate_artifacts(
    config: Config,
    start_date: str,
    end_date: str,
) -> AggregateArtifacts:
    strength = read_parquet(
        windowed_output_path(
            config.paths.aggregates_dir,
            config.paths.strength_file,
            start_date,
            end_date,
        )
    )
    synergy = read_parquet(
        windowed_output_path(
            config.paths.aggregates_dir,
            config.paths.synergy_file,
            start_date,
            end_date,
        )
    )
    counter = read_parquet(
        windowed_output_path(
            config.paths.aggregates_dir,
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
