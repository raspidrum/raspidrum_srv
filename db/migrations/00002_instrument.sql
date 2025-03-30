-- +goose Up
create table if not exists instrument (
  id          integer primary key autoincrement,
  uid         varchar(36) unique not null,
  key         varchar(128) not null,
  name        varchar(128) not null,
  fullname    varchar(512),
  type        varchar(16) not null,
  subtype     varchar(16) not null,
  midikey     varchar(16),
  description text,
  copyright   text,
  licence     text,
  credits     text,
  controls    text,
  layers      text
);

create table if not exists instrument_tag (
  id          integer primary key autoincrement,
  instrument  integer not null,
  name        varchar(16) not null,
  foreign key (instrument) references instrument(id) on delete cascade,
  unique (instrument, name)
);

create table if not exists kit_instrument (
  kit         integer not null,
  instrument  integer not null,
  foreign key (kit) references kit(id) on delete restrict,
  foreign key (instrument) references instrument(id) on delete cascade,
  unique (kit, instrument)
);

-- +goose Down
drop table kit_instrument;

drop table instrument_tag;

drop table instrument;