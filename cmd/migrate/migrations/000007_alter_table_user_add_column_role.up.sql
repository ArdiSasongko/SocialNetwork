alter table users
add column role varchar(255) references roles(name) default 'user';

update users
set role = (
    select name
    from roles
    where name = 'user'
);

alter table users
alter column role
drop default;

alter table users
alter column role 
set not null;