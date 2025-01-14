create table if not exists image_profile (
    user_id int not null,
    image_url varchar(255) not null,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    constraint fk_image_profile_user_id foreign key (user_id) references users(id) on delete cascade
);