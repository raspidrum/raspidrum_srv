-- +goose Up
/*
  iscustom = false - for kit, loaded with instrument samples. 
      Not custom kit always includes instruments with samples.
  iscustom = true - for kit, that uses preloaded kits (non custom) and instrument. 
      Custom kit doesn't have it's own instrument. 
      Custom kit is linked to preloaded instrument by preset.
*/
create table if not exists kit (
  id          integer primary key autoincrement,
  uid         varchar(36) unique not null,
  name        varchar(128) not null unique,
  iscustom    integer,
  description text,
  copyright   text,
  licence     text,
  credits     text,
  url         text
);


create table if not exists kit_tag (
  id          integer primary key autoincrement,
  kit         integer not null,
  name        varchar(16) not null,
  foreign key (kit) references kit(id) on delete cascade,
  unique (kit, name)
);


-- +goose Down
drop table kit_tag;

drop table kit;