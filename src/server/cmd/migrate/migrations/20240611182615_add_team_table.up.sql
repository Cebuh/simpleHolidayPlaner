CREATE TABLE IF NOT EXISTS teams (
    id UUID NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    createdAt TIMESTAMP not null DEFAULT UTC_TIMESTAMP
);
