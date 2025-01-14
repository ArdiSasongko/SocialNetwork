alter table posts
add column is_edited boolean not null default false;

update posts
set is_edited = false
where is_edited is null;