CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(100) NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS exercises (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT DEFAULT '',
    muscle_group VARCHAR(50) DEFAULT '',
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workouts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(200) NOT NULL,
    notes      TEXT DEFAULT '',
    date       DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sets (
    id          BIGSERIAL PRIMARY KEY,
    workout_id  BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id),
    reps        INT NOT NULL CHECK (reps > 0),
    weight_kg   NUMERIC(6,2) NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workouts_user_id  ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_sets_workout_id   ON sets(workout_id);
CREATE INDEX IF NOT EXISTS idx_exercises_user_id ON exercises(user_id);