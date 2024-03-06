CREATE TABLE Record(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    chunk VARCHAR(200),
    file_name VARCHAR(100),
    version INTEGER
);
