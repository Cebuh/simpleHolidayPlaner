CREATE TABLE IF NOT EXISTS users_teams (
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    addedAt TIMESTAMP not null DEFAULT UTC_TIMESTAMP,
    roletype int NOT NULL DEFAULT 1,
    CONSTRAINT team_user_user foreign key (user_id) references users(id),
    CONSTRAINT team_user_team foreign key (team_id) references teams(id),
    CONSTRAINT team_users_unique UNIQUE (team_id, user_id)
);
