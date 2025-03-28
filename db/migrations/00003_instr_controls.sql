-- +goose Up
create table if not exists instrument_control (
  id          integer primary key autoincrement,
  instrument  integer not null,
  name        varchar(128) not null,
  type        varchar(16),
  key         varchar(16) not null,
  foreign key (instrument) references instrument(id) on delete cascade,
  unique (instrument, name),
  unique (instrument, key)
);


-- +goose Down
drop table instrument_control;