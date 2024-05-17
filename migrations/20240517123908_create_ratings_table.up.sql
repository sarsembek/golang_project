CREATE TABLE IF NOT EXISTS ratings (
    car_id INT NOT NULL,
    stars INT NOT NULL,
    user_id INT NOT NULL,
    comment TEXT,
    CONSTRAINT fk_car FOREIGN KEY (car_id) REFERENCES cars (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (car_id, user_id)  -- Assuming combination of car_id and user_id is unique
);