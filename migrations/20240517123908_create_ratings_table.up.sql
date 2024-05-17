-- Migration to create ratings table without foreign key constraint
CREATE TABLE IF NOT EXISTS ratings (
    car_id INT NOT NULL,
    stars INT NOT NULL,
    user_id INT NOT NULL,
    comment TEXT,
    PRIMARY KEY (car_id, user_id)  -- Assuming combination of car_id and user_id is unique
);
