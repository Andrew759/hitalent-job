-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS departments
(
    id         bigserial PRIMARY KEY,
    name       VARCHAR(200) NOT NULL,
    parent_id  BIGINT
    CONSTRAINT fk_departments_departments
    REFERENCES departments(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_name_parent
    ON departments (name, parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_name_parent;
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd