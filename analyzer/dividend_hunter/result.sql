insert into analyzer_dividend_hunting_result(tag,
                                             figi,
                                             ticker,
                                             lastbuydate,
                                             recorddate,
                                             paymentdate,
                                             expectation)
values (:tag,
        :figi,
        :ticker,
        :lastbuydate,
        :recorddate,
        :paymentdate,
        :expectation)
on conflict (figi, tag) do nothing
;