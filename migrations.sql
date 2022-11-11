create table "user"(
    id serial primary key,
    name varchar(100)
);

create table balance(
    id serial primary key,
    user_id int not null,
    amount int not null,
    constraint user_fk foreign key(user_id) references "user"(id)
);

create table "transaction"(
    id serial primary key,
    user_id int not null,
    amount int not null,
    tx_type text check (tx_type in ('deposit', 'withdraw')),
    constraint user_fk foreign key(user_id) references "user"(id)
);

create table hold(
    id serial primary key,
    transaction_id int not null,
    balance_id int not null,
    amount int not null,
    tx_type text check (tx_type in ('deposit', 'withdraw')),
    constraint transaction_fk foreign key(transaction_id) references "transaction"(id),
    constraint balance_fk foreign key(balance_id) references balance(id)
);
