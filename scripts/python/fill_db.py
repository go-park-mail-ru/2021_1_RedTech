import psycopg2
from psycopg2 import sql
from faker import Faker
from random import random, randint, sample, choice, gauss
from datetime import date

from hashlib import sha256
import argparse


def parse_args():
    parser = argparse.ArgumentParser(description='Fill db tables')
    parser.add_argument('--users', type=int, dest='users',
                        help='count of new users')
    parser.add_argument('--movies', type=str, dest='movies',
                        help='path to file with movies data')
    parser.add_argument('--genres', type=str, dest='genres',
                        help='path to file with genres data')
    parser.add_argument('--actors', type=str, dest='actors',
                        help='path to file with actors data')
    parser.add_argument('--votes', dest='votes', action='store_false', 
                        help='flag if vote generation needed')
    parser.add_argument('--favs', dest='favs', action='store_false',
                        help='flag if users favs generation needed')
    parser.add_argument('--views', dest='views', action='store_false',
                        help='flag if movie views generation needed')

    return parser.parse_args()


def handle_args(args):
    conn = psycopg2.connect(dbname='netflix', user='redtech',
                            password='red_tech', host='localhost')
    c = conn.cursor()
    user_cnt = movie_cnt = actor_cnt = genre_cnt = 0

    if args.users is not None:
        user_cnt = create_users(c, args.users)
    if args.movies is not None:
        movie_cnt = create_movies(c, args.movies)
        create_movie_videos(c, movie_cnt)
    if args.genres is not None:
        genre_cnt = create_genres(c, args.genres)
        create_genre_links(c, movie_cnt, genre_cnt)
    if args.actors is not None:
        actor_cnt = create_actors(c, args.actors)
        create_actor_links(c, movie_cnt, actor_cnt)
    if args.votes:
        create_movie_votes(c, movie_cnt, user_cnt)
    if args.favs:
        create_user_favs(c, user_cnt, movie_cnt)
    if args.views:
        create_movie_views(c, user_cnt, movie_cnt)
        
    c.close()
    conn.commit()
    conn.close()
    return


def create_users(cursor, request_cnt):
    if request_cnt <= 0:
        return 0

    fake = Faker()
    result_cnt = 0
    for i in range(request_cnt):
        hasher = sha256()
        hasher.update(fake.word().encode('ascii'))
        username = fake.name()
        email = fake.email()
        password = hasher.digest()
        avatar = ''
        donate = False if randint(0, 4) else True
        try:
            cursor.execute("insert into users values(default, %s, %s, %s, %s, %s);",
                           [username, email, password, avatar, donate])
            result_cnt += 1
        except:
            print("it was an error while creating users")
            print("count: ", result_cnt)
            break

    print("Filling users table completed")
    return result_cnt


def create_movies(cursor, filepath):
    file = open(filepath, 'r')
    if file is None:
        return 0

    fake = Faker()
    result_cnt = 0
    for line in file:
        line = line.strip('\n').split(';')
        added = fake.date_between('-25y')
        title = line[0]
        descr = line[1]
        avatar = 'https://redioteka.com/static/media/img/movies/default.jpg'
        rating = random() * 10
        free = True if randint(0, 2) else False
        tip = choice([1, 2])
        year = int(line[2])
        dirs = line[3]
        cntrs = line[4]
        try:
            cursor.execute("insert into movies values(default, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);",
                           [title, descr, avatar, rating, free, tip, year, dirs, cntrs, added])
            result_cnt += 1
        except:
            print("it was an error while creating movies")
            print("count: ", result_cnt)
            break

    file.close()
    print("Filling movies table completed")
    return result_cnt


def create_movie_videos(cursor, request_cnt):
    for i in range(request_cnt):
        try:
            path = 'https://redioteka.com/static/media/movies/default.mp4'
            cursor.execute("insert into movie_videos values(default, %s, %s, %s);", [i + 1, path, int(gauss(6000, 4800))])
        except:
            print("it was an error while creating movie_videos")
            break
    print("Filling movie_videos table completed")
    return


def create_genres(cursor, filepath):
    file = open(filepath, 'r')
    if file is None:
        return 0

    result_cnt = 0
    for line in file:
        line = line.strip('\n').split(';')
        try:
            path = 'https://redioteka.com/static/media/img/genresAssets/' + line[2]
            cursor.execute("insert into genres values(default, %s, %s, %s);", [line[0], line[1], path])
            result_cnt += 1
        except:
            print("it was an error while creating genres")
            print("count: ", result_cnt)
            break

    file.close()
    print("Filling genres table completed")
    return result_cnt


def create_genre_links(cursor, movies, genres):
    movie_ids, movies = get_id_list(cursor, 'movies', movies)
    if movies < 10:
        print("Fill movies table firstly")
        return

    min_movie_cnt = movies // genres

    for i in range(genres):
        in_genre = randint(min_movie_cnt - 1, min_movie_cnt * 3)
        movies_sample = sample(movie_ids, in_genre)
        for m_id in movies_sample:
            try:
                cursor.execute("insert into movie_genres values(default, %s, %s);", [m_id, i + 1])
            except:
                print("it was an error while creating movie_genres")
                break

    print("Filling movie_genres table completed")
    return


def create_actors(cursor, filepath):
    file = open(filepath, 'r')
    if file is None:
        return 0

    result_cnt = 0
    for line in file:
        line = line.strip('\n').split(';')
        try:
            cursor.execute("insert into actors values(default, %s, %s);", line)
            result_cnt += 1
        except:
            print("it was an error while creating actors")
            print("count: ", result_cnt)
            break

    file.close()
    print("Filling actors table completed")
    return result_cnt


def create_actor_links(cursor, movies, actors):
    movie_ids, movies = get_id_list(cursor, 'movies', movies)
    if movies < 10:
        print("Fill movies table firstly")
        return
    min_movie_cnt = movies // actors

    for i in range(actors):
        actor_in = randint(min_movie_cnt - 1, min_movie_cnt * 3)
        movies_sample = sample(movie_ids, actor_in)
        for m_id in movies_sample:
            try:
                cursor.execute("insert into movie_actors values(default, %s, %s);", [m_id, i + 1])
            except:
                print("it was an error while creating movie_genres")
                break

    print("Filling movie_actors table completed")
    return


def create_movie_votes(cursor, movies, users):
    movie_ids, movies = get_id_list(cursor, 'movies', movies)
    user_ids, users = get_id_list(cursor, 'users', users)

    for m_id in movie_ids:
        votes_cnt = randint(0, users)
        users_voted = sample(user_ids, votes_cnt)
        for u_id in users_voted:
            vote = choice([-1, 1])
            try:
                cursor.execute("insert into movie_votes values(default, %s, %s, %s);", [u_id, m_id, vote])
            except:
                print("it was an error while creating movie_votes")
                break

    print("Filling movie_votes table completed")
    return


def create_user_favs(cursor, users, movies):
    movie_ids, movies = get_id_list(cursor, 'movies', movies)
    user_ids, users = get_id_list(cursor, 'users', users)

    for u_id in user_ids:
        favs_cnt = randint(0, movies // 2)
        movie_favs = sample(movie_ids, favs_cnt)
        for m_id in movie_favs:
            try:
                cursor.execute("insert into user_favs values(default, %s, %s);", [u_id, m_id])
            except:
                print("it was an error while creating user_favs")
                break

    print("Filling user_favs table completed")
    return

def create_movie_views(cursor, users, movies):
    movie_ids, movies = get_id_list(cursor, 'movies', movies)
    user_ids, users = get_id_list(cursor, 'users', users)

    for u_id in user_ids:
        views_cnt = randint(0, movies // 1.5)
        movie_views = sample(movie_ids, views_cnt)
        for m_id in movie_views:
            try:
                cursor.execute("insert into movie_views values(default, %s, %s);", [u_id, m_id])
            except:
                print("it was an error while creating movie_views")
                break
    
    print("Filling movie_views table completed")
    return

def get_id_list(cursor, table_name, count):
    cursor.execute(sql.SQL("select count(*) from {};").format(sql.Identifier(table_name)))
    count = max(count, cursor.fetchone()[0])
    id_list = [i + 1 for i in range(count)]
    return id_list, count


handle_args(parse_args())
