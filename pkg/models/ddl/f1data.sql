create table season
(
    id   int auto_increment
        primary key,
    year int not null
);

create table race
(
    id         int auto_increment
        primary key,
    season_id  int          not null,
    grand_prix varchar(255) not null,
    date       date null,
    constraint race_season_id_fk
        foreign key (season_id) references season (id)
);

create table race_result
(
    id            int auto_increment
        primary key,
    race_id       int           not null,
    position      varchar(10)   not null,
    driver_number int           not null,
    driver        varchar(255)  not null,
    team          varchar(255)  not null,
    laps          int           not null,
    time_retired  varchar(100)  not null,
    points        decimal(4, 2) not null,
    constraint race_result_race_id_fk
        foreign key (race_id) references race (id)
);

