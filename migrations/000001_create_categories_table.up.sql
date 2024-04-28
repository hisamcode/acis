CREATE TABLE IF NOT EXISTS categories (
    id smallint PRIMARY KEY,  
    title varchar(20) NOT NULL,
    emoji varchar(4) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);