-- Active: 1709698391238@@127.0.0.1@5432@filesync
CREATE TABLE Record(
    id INTEGER Primary Key,
    user_id INTEGER,
    chunk VARCHAR(200),
    file_name VARCHAR(100),
    version INTEGER
)