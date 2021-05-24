drop table if exists users cascade;
create table users
(
    id        serial not null primary key,
    username  varchar(64),
    email     varchar(64) unique,
    password  bytea,
    avatar    text,
    is_donate bool default false
);

drop table if exists movie_types cascade;
create table movie_types
(
    id   smallint not null primary key,
    type varchar(64)
);
insert into movie_types
values (1, 'Фильм'),
       (2, 'Сериал');

drop table if exists movies cascade;
create table movies
(
    id           serial not null primary key,
    title        text,
    description  text,
    avatar       text,
    rating       real,
    is_free      boolean,
    type         smallint,
    release_year smallint,
    directors    text,
    countries    text,
    add_date     date,
    constraint to_type foreign key (type) references movie_types (id) on delete set null
);

drop table if exists genres cascade;
create table genres
(
    id        serial not null primary key,
    name      varchar(64),
    label_rus varchar(64),
    image     text
);

drop table if exists movie_genres;
create table movie_genres
(
    id       serial not null primary key,
    movie_id int,
    genre_id int,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade,
    constraint to_genre foreign key (genre_id) references genres (id) on delete cascade
);

drop table if exists actors cascade;
create table actors
(
    id        serial not null primary key,
    firstname varchar(64),
    lastname  varchar(64),
    born      varchar(64),
    avatar    text
);

drop table if exists movie_actors;
create table movie_actors
(
    id       serial not null primary key,
    movie_id int,
    actor_id int,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade,
    constraint to_actor foreign key (actor_id) references actors (id) on delete cascade
);

drop table if exists user_favs;
create table user_favs
(
    id       serial not null primary key,
    user_id  int,
    movie_id int,
    constraint to_user foreign key (user_id) references users (id) on delete cascade,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade
);

drop table if exists movie_votes;
create table movie_votes
(
    id       serial not null primary key,
    user_id  int,
    movie_id int,
    value    smallint,
    constraint to_user foreign key (user_id) references users (id) on delete cascade,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade,
    constraint only_one unique (user_id, movie_id)
);

drop table if exists movie_videos;
create table movie_videos
(
    id       serial not null primary key,
    movie_id int,
    path     text,
    duration int,
    season   int,
    series   int,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade
);


drop table if exists movie_views;
create table movie_views
(
    id       serial not null primary key,
    user_id  int,
    movie_id int,
    constraint to_user foreign key (user_id) references users (id) on delete cascade,
    constraint to_movie foreign key (movie_id) references movies (id) on delete cascade,
    constraint only_one_view_per_user unique (user_id, movie_id)
);

drop table if exists subscriptions;
create table subscriptions
(
    id serial primary key,
    user_id int,
    expires int,
    constraint to_user foreign key (user_id) references users (id) on delete cascade
);
