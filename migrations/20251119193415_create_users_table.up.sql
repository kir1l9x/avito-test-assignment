CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY,
    username   TEXT NOT NULL,
    team_name    TEXT NOT NULL REFERENCES teams(name) ON DELETE RESTRICT,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE
);
