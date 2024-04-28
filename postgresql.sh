# !/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    \set acis_db `echo $POSTGRES_ACIS_DB`
    CREATE DATABASE :acis_db;
    \c :acis_db;
    \set acis_user `echo $POSTGRES_ACIS_USER`
    \set acis_password `echo $POSTGRES_ACIS_PASSWORD`
    CREATE ROLE :acis_user WITH LOGIN PASSWORD :'acis_password';
    create extension if not exists "citext";
EOSQL