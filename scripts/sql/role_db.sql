drop role if exists redtech;
create role redtech with password 'red_tech';
alter role redtech with login superuser;
DROP DATABASE netflix;
CREATE DATABASE netflix
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;
