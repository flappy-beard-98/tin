insert into collector_dividends(figi,
                                dividendnet,
                                paymentdate,
                                declareddate,
                                lastbuydate,
                                dividendtype,
                                recorddate,
                                regularity,
                                closeprice,
                                yieldvalue,
                                createdat)
values (:figi,
        :dividendnet,
        :paymentdate,
        :declareddate,
        :lastbuydate,
        :dividendtype,
        :recorddate,
        :regularity,
        :closeprice,
        :yieldvalue,
        :createdat)
on conflict (figi,createdat) do nothing
;