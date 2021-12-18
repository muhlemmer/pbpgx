create table simple_ro (
    id integer primary key not null,
    title text null,
    data text null
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
    data text null
);


create table products (
    id bigint primary key not null,
    title text not null,
    price double precision null
);

insert into products (id, title, price) values 
    (1, 'one', 9.99),
    (2, 'two', 10.45),
    (3, 'three', null),
    (4, 'four', 100),
    (5, 'five', 0.90);
