create table if not exists images_post (
    image_name varchar(255) not null,
    post_id int not null,
    image_url varchar(255) not null,
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone not null default now(),
    constraint fk_images_post_post_id foreign key (post_id) references posts(id) on delete cascade
);