-- migrate:up
CREATE TABLE IF NOT EXISTS car (
    id SERIAL PRIMARY KEY,
    brand VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    color VARCHAR(50) NOT NULL,
    body_style VARCHAR(100),
    engine_size DECIMAL(10,2),
    weight DECIMAL(10,2),
    base_price INT,
    fuel_capacity INT,
    horsepower INT,
    torque INT,
    acceleration INT,
    top_speed INT
);

