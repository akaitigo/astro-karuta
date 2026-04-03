CREATE TABLE IF NOT EXISTS cards (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    category    TEXT NOT NULL CHECK (category IN ('constellation', 'planet', 'phenomenon')),
    reading_text TEXT NOT NULL,
    image_url   TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    magnitude   DOUBLE PRECISION,
    distance    TEXT,
    best_season TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS decks (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    card_ids   UUID[] NOT NULL DEFAULT '{}',
    seasonal   BOOLEAN NOT NULL DEFAULT false,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    valid_to   TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '1 year',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS games (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_code      TEXT NOT NULL UNIQUE,
    status         TEXT NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'playing', 'finished')),
    deck_id        UUID NOT NULL REFERENCES decks(id),
    current_card_id UUID,
    remaining_ids  UUID[] NOT NULL DEFAULT '{}',
    time_limit_sec INT NOT NULL DEFAULT 300,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at     TIMESTAMPTZ,
    finished_at    TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS game_players (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id      UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL REFERENCES users(id),
    score        INT NOT NULL DEFAULT 0,
    captured_ids UUID[] NOT NULL DEFAULT '{}',
    is_connected BOOLEAN NOT NULL DEFAULT true,
    disconnect_at TIMESTAMPTZ,
    UNIQUE(game_id, user_id)
);

CREATE TABLE IF NOT EXISTS collections (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    card_id     UUID NOT NULL REFERENCES cards(id),
    obtained_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    source      TEXT NOT NULL DEFAULT 'game' CHECK (source IN ('game', 'mission')),
    UNIQUE(user_id, card_id)
);

CREATE TABLE IF NOT EXISTS observation_missions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id    UUID NOT NULL REFERENCES cards(id),
    title      TEXT NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to   TIMESTAMPTZ NOT NULL,
    latitude   DOUBLE PRECISION NOT NULL DEFAULT 0,
    longitude  DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_collections_user_id ON collections(user_id);
CREATE INDEX idx_game_players_game_id ON game_players(game_id);
CREATE INDEX idx_games_room_code ON games(room_code);
CREATE INDEX idx_observation_missions_valid ON observation_missions(valid_from, valid_to);
