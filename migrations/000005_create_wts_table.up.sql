CREATE TABLE IF NOT EXISTS wts (
    id smallint PRIMARY KEY,
    name text not null,
    detail text not null
);

insert into wts (id, name, detail) values (1, 'in', 'kredit');
insert into wts (id, name, detail) values (2, 'out', 'debit');

