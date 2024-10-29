create table season
(
    id         int      not null auto_increment,
    year       int      not null,
    updated_at datetime not null default now(),
    primary key (id)
);

