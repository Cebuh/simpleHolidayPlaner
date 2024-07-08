CREATE TABLE IF NOT EXISTS vacation_requests (
    id UUID NOT NULL PRIMARY KEY,
    requestedFrom UUID NOT NULL,
    toUserId UUID NOT NULL,
    teamId UUID NOT NULL,
    fromDate TIMESTAMP not null,
    toDate TIMESTAMP not null,
    info varchar(255), 
    requestStatus int not null default 0,
    changedAt TIMESTAMP,
    createdAt TIMESTAMP not null DEFAULT UTC_TIMESTAMP,
    CONSTRAINT requests_from foreign key (requestedFrom) references users(id),
    CONSTRAINT requests_teamId foreign key (teamId) references teams(id),
    CONSTRAINT requests_unique UNIQUE (requestedFrom, teamId)
);
