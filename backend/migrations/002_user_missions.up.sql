CREATE TABLE IF NOT EXISTS user_missions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    mission_id   TEXT NOT NULL,
    card_id      UUID NOT NULL REFERENCES cards(id),
    title        TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'expired')),
    valid_from   TIMESTAMPTZ NOT NULL,
    valid_to     TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_missions_user_id ON user_missions(user_id);
CREATE INDEX idx_user_missions_status ON user_missions(user_id, status);
CREATE INDEX idx_user_missions_valid ON user_missions(valid_from, valid_to);
