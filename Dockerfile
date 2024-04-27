FROM postgres:16.2-alpine

COPY postgresql.sh /docker-entrypoint-initdb.d/

RUN chmod 755 /docker-entrypoint-initdb.d/postgresql.sh