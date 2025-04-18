-- +goose Up
create table if not exists kit_preset (
  id          integer primary key autoincrement,
  uid         varchar(36) unique not null,
  kit         integer not null,
  name        varchar(128) not null,
  foreign key (kit) references kit(id) on delete cascade,
  unique (kit, name),
  unique (kit, uid)
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

create view v_kit_preset as
select p.*, k.uid as kit_uid, k.name as kit_name, k.iscustom as kit_iscustom
  from kit_preset p
  join kit k on p.kit = p.id;

create view v_preset_instrument as
select pi.*, chn.key as channel_key, 
			 i.uid as instrument_uid, i.key as instrument_key, i.name as instrument_name,
       i.midikey as instrument_midikey,
			 i.controls as instrument_controls, i.layers as instrument_layers
	from preset_instrument pi
	join preset_channel chn on chn.preset = pi.preset and chn.id = pi.channel
	join instrument i on i.id = pi.instrument;

-- +goose Down
drop view v_kit_preset;

drop view v_preset_instrument;

drop table preset_instrument;

drop table preset_channel;

drop table kit_preset;