create table if not exists roles(
    id serial primary key,
    name varchar(255) not null unique,
    level int not null,
    description text not null
);

insert into roles(name, level, description)
values ('user', 1, 'a user can only create, delete, update his own data and read all data');

insert into roles(name, level, description)
values ('moderator', 2, 'a moderator can only update other users data and read all data');

insert into roles(name, level, description)
values ('admin', 3, 'an admin can do anything');