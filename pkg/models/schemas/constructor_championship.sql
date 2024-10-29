create table constructor_championship
(
    id         int           not null auto_increment,
    season_id  int           not null,
    position   int           not null,
    name       varchar(255)  not null,
    points     decimal(6, 2) not null,
    updated_at datetime      not null default now(),
    primary key (id),
    constraint team_championship_season_id_fk
        foreign key (season_id) references season (id)
);

