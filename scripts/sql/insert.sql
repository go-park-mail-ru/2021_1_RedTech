insert into users values(default, 'testuser', 'test@mail.ru', 'pass', '', default),
						 (default, 'superuser', 'mail@gmail.com', 'hash from pass', '', true);
insert into movies values(default, 'Заголовок фильма', 'какое-то описание', '', 0, true, 1, 2020, 'Кто-то', 'Россия', '2021-03-31'),
						 (default, 'Тест-сериал', 'Что это?', '', 5, false, 2, 2021, 'RedTech', 'Россия, Татарстан', '2021-04-01'),
						 (default, 'Умерь свой энтузиазм', 'мем', '', 10, true, 2, 2000, 'Robert B. Weide', 'США', '2007-09-24');
insert into genres values(default, 'Horror', 'Ужасы', 'https://redioteka.com/static/media/img/genres/horror.svg'), (default, 'Comedy', 'Комедии', 'https://redioteka.com/static/media/img/genres/comedy.svg'), (default, 'Sci-fi', 'Фантастические', 'https://redioteka.com/static/media/img/genres/sci-fi.svg');
insert into movie_genres values(default, 1, 1), (default, 2, 1), (default, 2, 3), (default, 3, 2);
insert into actors values(default, 'Джон', 'Сина'), (default, 'Ларри', 'Дэвид'), (default, 'Шерил', 'Хайнс');
insert into movie_actors values(default, 1, 1), (default, 3, 2), (default, 3, 3);
insert into user_favs values(default, 1, 2), (default, 2, 1), (default, 2, 2), (default, 2, 3);
insert into movie_votes values(default, 1, 1, -1), (default, 1, 2, 1), (default, 1, 3, -1), 
							(default, 2, 2, 1), (default, 2, 3, 1);
insert into movie_videos values(default, 1, '', 600);
select * from users;
select * from movie_types;
select * from movies;
select * from genres;
select * from movie_genres;
select * from actors;
select * from movie_actors;
select * from user_favs;
select * from movie_votes;
select * from movie_videos;