create table if not exists user_activities(
    id serial primary key,
    user_id int not null,
    post_id int not null,
    is_liked boolean,
    is_disliked boolean,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    constraint fk_user_activities_user_id foreign key (user_id) references users(id),
    constraint fk_user_activities_post_id foreign key (post_id) references posts(id)
);