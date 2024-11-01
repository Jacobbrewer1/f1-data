create table season
(
    id   int auto_increment
        primary key,
    year int not null
);

create table constructor_championship
(
    id        int auto_increment
        primary key,
    season_id int           not null with system versioning,
    position  int           not null with system versioning,
    name      varchar(255)  not null with system versioning,
    points    decimal(6, 2) not null with system versioning,
    constraint constructor_championship_season_id_fk
        foreign key (season_id) references season (id)
);

create table driver_championship
(
    id          int auto_increment
        primary key,
    season_id   int           not null with system versioning,
    position    int           not null with system versioning,
    driver      varchar(255)  not null with system versioning,
    driver_tag  varchar(3)    not null with system versioning,
    nationality varchar(3)    not null with system versioning,
    team        varchar(255)  not null with system versioning,
    points      decimal(6, 2) not null with system versioning,
    constraint driver_championship_season_id_fk
        foreign key (season_id) references season (id)
);

create table race
(
    id         int auto_increment
        primary key,
    season_id  int          not null with system versioning,
    grand_prix varchar(255) not null with system versioning,
    date       date         not null with system versioning,
    constraint race_season_id_fk
        foreign key (season_id) references season (id)
);

create table race_result
(
    id            int auto_increment
        primary key,
    race_id       int           not null with system versioning,
    position      varchar(10)   not null with system versioning,
    driver_number int           not null with system versioning,
    driver        varchar(255)  not null with system versioning,
    driver_tag    varchar(3)    not null with system versioning,
    team          varchar(255)  not null with system versioning,
    laps          int           not null with system versioning,
    time_retired  varchar(100)  not null with system versioning,
    points        decimal(4, 2) not null with system versioning,
    constraint race_result_race_id_fk
        foreign key (race_id) references race (id)
);

