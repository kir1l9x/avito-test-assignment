CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id    TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id  TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pr_id, user_id)
);
