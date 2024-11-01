alter table race_result drop system versioning;
alter table race_result drop column updated_at;
alter table race_result add system versioning;

alter table race drop system versioning;
alter table race drop column updated_at;
alter table race add system versioning;

alter table driver_championship drop system versioning;
alter table driver_championship drop column updated_at;
alter table driver_championship add system versioning;

alter table constructor_championship drop system versioning;
alter table constructor_championship drop column updated_at;
alter table constructor_championship add system versioning;
