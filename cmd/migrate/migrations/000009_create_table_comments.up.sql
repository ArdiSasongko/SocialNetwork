create table if not exists comments(
    id serial primary key,
    user_id int not null,
    post_id int not null,
    is_edited boolean not null default false,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    constraint fk_comments_user_id foreign key(user_id) references users(id) on delete cascade,
    constraint fk_comments_post_id foreign key(post_id) references posts(id) on delete cascade
);