create table if not exists posts (
    id serial primary key,
    user_id int not null,
    title varchar(255) not null,
    content text not null,
    tags varchar(255)[],
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    constraint fk_posts_user_id foreign key (user_id) references users(id) on delete cascade
);