insert into collector_last_prices(figi,
                                  price,
                                  time,
                                  instrumentuid)
values (:figi,
        :price,
        :time,
        :instrumentuid)
on conflict (figi) do nothing
;