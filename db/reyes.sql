create table toys
(
    toy_id     integer not null
        constraint toys_pk
            primary key autoincrement,
    toy_name   TEXT    not null,
    age_min    integer not null,
    age_max    integer not null,
    image1     TEXT,
    image2     TEXT,
    source_url TEXT,
    deleted integer default 0
);

create table volunteers
(
    volunteer_id integer not null
        constraint volunteers_pk
            primary key autoincrement,
    name         TEXT    not null,
    email        TEXT    not null,
    phone        TEXT,
    address      TEXT    not null,
    address2     TEXT,
    country      TEXT    not null,
    state        TEXT    not null,
    city         TEXT    not null,
    province     TEXT,
    zip_code     TEXT    not null,
    secret       TEXT    not null,
    passcode     integer not null,
    deleted integer default 0
);

create table codes
(
    code_id    integer           not null
        constraint codes_pk
            primary key autoincrement,
    code       TEXT              not null,
    expiration TEXT              not null,
    used       integer default 0 not null,
    cancelled  integer default 0,
    deleted integer default 0
);

create unique index codes_code_uindex
    on codes (code);

create table volunteer_codes
(
    volunteer_code_id integer           not null
        constraint volunteer_codes_pk
            primary key autoincrement,
    volunteer_id      integer           not null
        constraint volunteer_codes_volunteers_volunteer_id_fk
            references volunteers,
    code_id           integer           not null
        constraint volunteer_codes_codes_code_id_fk
            references codes,
    deleted           integer default 0 not null
);

create table orders
(
    order_id     integer not null
        constraint orders_pk
            primary key autoincrement,
    toy_id       integer not null
        constraint orders_toys_fk
            references toys,
    volunteer_id integer not null
        constraint orders_volunteers_fk
            references volunteers,
    code_id      integer not null
        constraint orders_codes_fk
            references codes,
    order_date   TEXT    not null,
    shipped      integer default 0 not null,
    shipped_date TEXT,
    deleted integer default 0
);
