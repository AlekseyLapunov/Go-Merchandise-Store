create table employees (
    id       serial primary key,
    login    varchar(25) unique not null,
    password varchar(62) not null,
    coins    int default 0
);

create table merch (
    id   serial primary key,
    name varchar(20) unique not null,
    cost int check (cost >= 0)
);

create table purchases (
    id       serial primary key,
    emp_id   int references employees(id),
    merch_id int references merch(id)
);

create table operations (
    id          serial primary key,
    recv_emp_id int references employees(id),
    send_emp_id int references employees(id),
    amount      int check(amount > 0)
);
