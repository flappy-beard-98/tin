create table if not exists collector_last_prices
(
    figi          text primary key,
    price         real,
    time          text,
    instrumentuid text
) without rowid;
