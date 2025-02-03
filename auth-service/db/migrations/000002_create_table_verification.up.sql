CREATE TABLE IF NOT EXISTS user_verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expire_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_fk_user_verification_user_id ON user_verifications (user_id);
