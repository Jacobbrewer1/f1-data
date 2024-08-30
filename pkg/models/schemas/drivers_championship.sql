create table driver_championship
(
    id          int           not null auto_increment,
    season_id   int           not null,
    driver      varchar(255)  not null,
    driver_tag  varchar(3)    not null,
    nationality varchar(3)    not null,
    team        varchar(255)  not null,
    points      decimal(6, 2) not null,
    primary key (id),
    constraint driver_championship_season_id_fk
        foreign key (season_id) references season (id)
);

