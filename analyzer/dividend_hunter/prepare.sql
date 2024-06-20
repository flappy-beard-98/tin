with next_dividends as (select d.figi,
                               d.dividendnet,
                               d.lastbuydate,
                               d.recorddate,
                               d.paymentdate,
                               row_number() over (partition by d.figi order by d.lastbuydate ) as rownum
                        from collector_dividends as d
                        where lastbuydate >= date('now'))
insert
into analyzer_dividend_hunting_base(figi,
                                    dividendnet,
                                    lastbuydate,
                                    recorddate,
                                    paymentdate,
                                    price,
                                    lot,
                                    unitprice,
                                    ticker,
                                    name)
select d.figi,
       d.dividendnet,
       d.lastbuydate,
       d.recorddate,
       d.paymentdate,
       p.price,
       s.lot,
       p.price * s.lot as unitprice,
       s.ticker,
       s.name
from next_dividends d
         join collector_last_prices p on p.figi = d.figi
         join collector_shares s on s.figi = d.figi
where d.rownum = 1
on conflict (figi) do nothing;
