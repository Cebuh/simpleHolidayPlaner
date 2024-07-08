CREATE TABLE IF NOT EXISTS vacation_approvals (
    request_id UUID NOT NULL,
    approver_id UUID NOT NULL,
    status int not null default 0,
    changedAt TIMESTAMP,
    CONSTRAINT request_approval foreign key (request_id) references vacation_requests(id),
    CONSTRAINT approvals_users foreign key (approver_id) references users(id),
    CONSTRAINT approvals_requests_unique UNIQUE (request_id, approver_id)
);
