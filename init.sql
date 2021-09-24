CREATE TABLE IF NOT EXISTS Links (
    short_url varchar(11) NOT NULL PRIMARY KEY,
    url text NOT NULL DEFAULT ''
);

