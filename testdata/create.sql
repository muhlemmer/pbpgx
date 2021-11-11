create table simple_ro (
    id integer primary key not null,
    title text null
);

insert into simple_ro (id, title) values 
    (1, 'one'),
    (2, 'two'),
    (3, null),
    (4, 'four'),
    (5, 'five');

create table simple_rw (
    id integer primary key not null,
    title text null
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
