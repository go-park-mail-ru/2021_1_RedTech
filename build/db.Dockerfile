FROM ubuntu:20.04
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV DEBIAN_FRONTEND=noninteractive
ENV PGVER 12
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB netflix
ENV POSTGRES_USER postgres

RUN apt-get update && apt-get install -y postgresql-$PGVER python3 python3-pip
RUN pip3 install psycopg2-binary faker
EXPOSE $POSTGRES_PORT

USER $POSTGRES_USER

COPY scripts/sql/role_db.sql role_db.sql
COPY scripts/sql/tables.sql tables.sql
COPY scripts/python/ .

RUN service postgresql start &&\
    psql -U $POSTGRES_USER -f role_db.sql &&\
    psql -U $POSTGRES_USER -d $POSTGRES_DB -a -f tables.sql &&\
    python3 fill_db.py --users 100 --movies movies.txt --actors actors.txt --genres genres.txt --views --votes --favs

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "shared_buffers=256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "full_page_writes=off" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD /usr/lib/postgresql/$PGVER/bin/postgres -D /var/lib/postgresql/$PGVER/main -c config_file=/etc/postgresql/$PGVER/main/postgresql.conf
