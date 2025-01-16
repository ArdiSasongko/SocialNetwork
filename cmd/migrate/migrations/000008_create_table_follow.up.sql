create table if not exists follows(
    user_id int not null,
    follower_id int not null,
    created_at timestamp(0) with time zone not null default now(),
    primary key(user_id, follower_id)
);