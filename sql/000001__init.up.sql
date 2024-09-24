create table if not exists runeprice
(
    id        serial primary key,
    rune_name varchar(255) not null,
    server    varchar(255) not null,
    date      timestamp    not null default current_timestamp,
    price     int          not null,
    constraint unique_rune_server_date unique (rune_name, server, date)
);