-- migrate:up
CREATE TABLE IF NOT EXISTS car_history (
    id SERIAL PRIMARY KEY,
    car_id INT NOT NULL,
    date TIMESTAMP NOT NULL,
    type VARCHAR(20) NOT NULL,
    details TEXT,
    service_type VARCHAR(50),
    service_cost DECIMAL(10,2),
    service_notes TEXT
);