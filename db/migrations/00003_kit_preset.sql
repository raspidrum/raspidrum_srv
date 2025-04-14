-- +goose Up
create table if not exists kit_preset (
  id          integer primary key autoincrement,
  uid         varchar(36) unique not null,
  kit         integer not null,
  name        varchar(128) not null,
  foreign key (kit) references kit(id) on delete cascade,
  unique (kit, name)
);

create table if not exists preset_channel (
  id          integer primary key autoincrement,
  preset      integer not null,
  key         varchar(16) not null,
  name        varchar(16) not null,
  controls    text not null,
  foreign key (preset) references kit_preset(id) on delete cascade,
  unique (preset, key)
);

create table if not exists preset_instrument (
  id          integer primary key autoincrement,
  preset      integer not null,
  channel     integer not null,
  instrument  integer not null,
  name        varchar(16) not null,
  midikey     varchar(16),
  controls    text not null,
  layers      text,
  foreign key (preset) references kit_preset(id) on delete cascade,
  foreign key (instrument) references instrument(id) on delete restrict,
  foreign key (channel) references preset_channel(id) on delete cascade,
  unique (preset, name)
);


-- +goose Down
drop table preset_instrument;

drop table preset_channel;

drop table kit_preset;