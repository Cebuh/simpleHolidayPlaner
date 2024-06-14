CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) not null ,
    password varchar(255) not null,
    createdAt TIMESTAMP not null DEFAULT current_timestamp
);
