create table race_result
(
    id            int           not null auto_increment,
    race_id       int           not null,
    position      varchar(10)   not null,
    driver_number int           not null,
    driver        varchar(255)  not null,
    driver_tag    varchar(3)    not null,
    team          varchar(255)  not null,
    laps          int           not null,
    time_retired  varchar(100)  not null,
    points        decimal(4, 2) not null,
    primary key (id),
    constraint race_result_race_id_fk
        foreign key (race_id) references race (id)
);

