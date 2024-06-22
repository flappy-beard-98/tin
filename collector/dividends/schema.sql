create table if not exists collector_dividends
(
    figi         text,
    dividendnet  real,
    paymentdate  text,
    declareddate text,
    lastbuydate  text,
    dividendtype text,
    recorddate   text,
    regularity   text,
    closeprice   real,
    yieldvalue   real,
    createdat    text,

    constraint collector_dividends_pk
        primary key (figi, createdat)
) without rowid;