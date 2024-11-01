alter table race_result drop system versioning;
alter table race_result add column updated_at datetime not null default now();
alter table race_result add system versioning;

alter table race drop system versioning;
alter table race add column updated_at datetime not null default now();
alter table race add system versioning;

alter table driver_championship drop system versioning;
alter table driver_championship add column updated_at datetime not null default now();
alter table driver_championship add system versioning;

alter table constructor_championship drop system versioning;
alter table constructor_championship add column updated_at datetime not null default now();
alter table constructor_championship add system versioning;
