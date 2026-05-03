from datetime import date, timedelta


def target_dates(
    lookback_days: int,
    end_date: date | None = None,
    include_today: bool = True,
) -> list[date]:
    if end_date is None:
        end_date = date.today()
        if not include_today:
            end_date = end_date - timedelta(days=1)

    start_date = end_date - timedelta(days=lookback_days - 1)
    days = []
    current = start_date
    while current <= end_date:
        days.append(current)
        current += timedelta(days=1)
    return days


def prior_window_for_day(target_day: date, prior_days: int) -> tuple[date, date]:
    start_date = target_day - timedelta(days=prior_days)
    end_date = target_day - timedelta(days=1)
    return start_date, end_date
