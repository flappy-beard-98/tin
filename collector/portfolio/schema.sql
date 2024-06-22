create table if not exists collector_portfolios
(
    totalamountshares     REAL,
    totalamountbonds      REAL,
    totalamountetf        REAL,
    totalamountcurrencies REAL,
    totalamountfutures    REAL,
    expectedyield         REAL,
    accountid             TEXT PRIMARY KEY,
    totalamountoptions    REAL,
    totalamountsp         REAL,
    totalamountportfolio  REAL
) without rowid;

create table if not exists collector_virtual_positions
(
    accountid                TEXT,
    positionuid              TEXT PRIMARY KEY,
    instrumentuid            TEXT,
    figi                     TEXT,
    instrumenttype           TEXT,
    quantity                 REAL,
    averagepositionprice     REAL,
    expectedyield            REAL,
    expectedyieldfifo        REAL,
    expiredate               TEXT,
    currentprice             REAL,
    averagepositionpricefifo REAL
) without rowid;

create table if not exists collector_positions
(
    accountid                TEXT,
    figi                     TEXT,
    instrumenttype           TEXT,
    quantity                 REAL,
    averagepositionprice     REAL,
    expectedyield            REAL,
    currentnkd               REAL,
    currentprice             REAL,
    averagepositionpricefifo REAL,
    blocked                  BOOLEAN,
    blockedlots              REAL,
    positionuid              TEXT primary key,
    instrumentuid            TEXT,
    varmargin                REAL,
    expectedyieldfifo        REAL
) without rowid;

