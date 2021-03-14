create table users (
    id serial primary key,
    username varchar(20) not null unique,
    password text not null
);

-- seeders
insert into users (username, password) 
    values ('test@example.com', '1234567890');
insert into users (username, password) 
    values ('test123@example.com', '1234567890');