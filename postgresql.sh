# !/bin/bash
set -e
# \c :acis_db;

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    \set acis_db `echo $POSTGRES_ACIS_DB`
    CREATE DATABASE :acis_db;
    \set acis_user `echo $POSTGRES_ACIS_USER`
    \set acis_password `echo $POSTGRES_ACIS_PASSWORD`
    CREATE USER :acis_user WITH PASSWORD :'acis_password';
    GRANT ALL PRIVILEGES ON DATABASE :acis_db TO :acis_user;
    ALTER DATABASE :acis_db OWNER to :acis_user;
    create extension if not exists "citext";
EOSQL