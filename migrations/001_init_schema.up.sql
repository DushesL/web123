-- migrations/001_init_schema.up.sql

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    role          VARCHAR(50) NOT NULL CHECK (role IN ('client', 'psychologist')),
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS slots (
    id              SERIAL PRIMARY KEY,
    psychologist_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_at        TIMESTAMP NOT NULL,
    end_at          TIMESTAMP NOT NULL,
    is_booked       BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT unique_psychologist_slot UNIQUE (psychologist_id, start_at),
    CONSTRAINT check_slot_range CHECK (end_at > start_at)
);

CREATE INDEX idx_slots_psychologist_id ON slots (psychologist_id);
CREATE INDEX idx_slots_start_at ON slots (start_at);

CREATE TABLE IF NOT EXISTS appointments (
    id              SERIAL PRIMARY KEY,
    client_id       INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slot_id         INT NOT NULL REFERENCES slots(id) ON DELETE RESTRICT,
    status          VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled', 'rescheduled')),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_appointments_client_id ON appointments (client_id);
CREATE INDEX idx_appointments_slot_id ON appointments (slot_id);