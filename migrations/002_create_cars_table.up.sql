CREATE TABLE IF NOT EXISTS cars (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    license_plate VARCHAR(20) NOT NULL,
    color VARCHAR(50),
    size VARCHAR(50),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(tg_user_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_cars_user_id ON cars(user_id);
CREATE INDEX idx_cars_license_plate ON cars(license_plate);