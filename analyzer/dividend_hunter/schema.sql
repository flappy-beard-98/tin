create table if not exists analyzer_dividend_hunting_base
(
    figi        text primary key,
    dividendnet real,
    lastbuydate text,
    recorddate  text,
    paymentdate text,
    price       real,
    lot         int,
    unitprice   real,
    ticker      text,
    name        text
) without rowid;

create table if not exists analyzer_dividend_hunting_result
(
    tag         text,
    figi        text,
    ticker      text,
    lastbuydate text,
    recorddate  text,
    paymentdate text,
    expectation real,

    primary key (figi, tag)
) without rowid;
