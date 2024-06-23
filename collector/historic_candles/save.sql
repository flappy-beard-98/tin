insert into collector_historic_candles (uid,
                                        open,
                                        high,
                                        low,
                                        close,
                                        volume,
                                        time,
                                        iscomplete,
                                        candlesource)
values (:uid,
        :open,
        :high,
        :low,
        :close,
        :volume,
        :time,
        :iscomplete,
        :candlesource)
on conflict (uid, time) do nothing;
