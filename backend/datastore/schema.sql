CREATE SCHEMA something;

CREATE TABLE something.somethings (
    id BIGSERIAL PRIMARY KEY,
    external_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    text TEXT NOT NULL
);

CREATE UNIQUE INDEX unique_something ON something.somethings (external_id);