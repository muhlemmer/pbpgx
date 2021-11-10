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
