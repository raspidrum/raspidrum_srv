
create table if not exists kit (
  id          integer primary key autoincrement,
  uid         varchar(36) unique not null,
  name        varchar(128) not null,
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
  name        varchar(16) not null
  foreign key (kit) references kit(id) on delete cascade
);
