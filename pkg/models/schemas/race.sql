create table race
(
    id         int          not null auto_increment,
    season_id  int          not null,
    grand_prix varchar(255) not null,
    date       date         not null,
    primary key (id),
    constraint race_season_id_fk
        foreign key (season_id) references season (id)
);

