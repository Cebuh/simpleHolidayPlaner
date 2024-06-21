CREATE TABLE IF NOT EXISTS invites (
    id UUID NOT NULL PRIMARY KEY,
    fromUserId UUID NOT NULL,
    toUserId UUID NOT NULL,
    teamId UUID NOT NULL,
    status int NOT NULL DEFAULT (0),
    inviteType int NOT NULL, 
    createdAt TIMESTAMP not null DEFAULT UTC_TIMESTAMP,
    changedAt TIMESTAMP,
    CONSTRAINT invites_user_from foreign key (fromUserId) references users(id),
    CONSTRAINT invites_user_to foreign key (toUserId) references users(id),
    CONSTRAINT invites_teamId foreign key (teamId) references teams(id),
    CONSTRAINT invites_unique UNIQUE (fromUserId, toUserId, teamId)
);
