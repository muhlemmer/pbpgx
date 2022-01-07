create table simple_ro (
    id integer primary key not null,
    title text null,
    data text null,
    created timestamptz null default now()
);

insert into simple_ro (id, title, data) values 
    (1, 'one', 'foo bar'),
    (2, 'two', null),
    (3, null, 'golden triangle'),
    (4, 'four', 'hello world'),
    (5, 'five', 'five is a four letter word');

create table simple_rw (
    id integer primary key not null,
    title text null,
    data text null,
    created timestamptz null default now()
);


create table products (
    id bigserial primary key not null,
    title text not null,
    price double precision null,
    created timestamptz null default now()
);

insert into products (title, price, created) values 
    ('one', 9.99, '2022-01-07 13:47:07'),
    ('two', 10.45, '2022-01-07 13:47:08'),
    ('three', null, '2022-01-07 13:47:09'),
    ('four', 100, '2022-01-07 13:47:10'),
    ('five', 0.90,'2022-01-07 13:47:11');

create table unsupported (
    bl boolean[],
    i32 integer[],
    i64 bigint[],
    f float[],
    d real[],
    s text[],
    bt bytea[],
    u32 integer[],
    u64 bigint[],
    sup bytea,
    ts timetz
);