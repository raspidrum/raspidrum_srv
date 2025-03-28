-- +goose Up
create table if not exists instrument_layer (
  id          integer primary key autoincrement,
  instrument  integer not null,
  name        varchar(128) not null,
  midikey     varchar(16),
  foreign key (instrument) references instrument(id) on delete cascade,
  unique (instrument, name),
);

create table if not exists layer_control (
  id          integer primary key autoincrement,
  layer       integer not null,
  name        varchar(128) not null,
  type        varchar(16),
  key         varchar(16) not null,
  foreign key (layer) references instrument_layer(id) on delete cascade,
  unique (layer, name),
  unique (layer, key)
);

-- +goose Down
drop table instrument_layer;