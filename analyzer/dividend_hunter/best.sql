;with totals as (
    select
        SUM(expectation) as total,
        tag
    from analyzer_dividend_hunting_result
    group by tag
)
 select
     tag,
     figi,
     ticker,
     lastbuydate,
     recorddate,
     paymentdate,
     cast(expectation as int) as expectation
 from analyzer_dividend_hunting_result
 where tag = (select tag from totals order by total desc limit 1)
 order by lastbuydate