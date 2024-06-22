create table if not exists collector_accounts
(
    id          text primary key,
    type        text,
    name        text,
    status      text,
    openeddate  text,
    closeddate  text,
    accesslevel text
) without rowid;