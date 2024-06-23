create table if not exists collector_historic_candles
(
    uid          text,
    open         real,
    high         real,
    low          real,
    close        real,
    volume       integer,
    time         text,
    iscomplete   boolean,
    candlesource text,

    primary key (uid, time)
) without rowid;