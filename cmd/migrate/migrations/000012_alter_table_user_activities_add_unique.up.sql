alter table user_activities
add constraint unique_user_post unique (user_id, post_id);