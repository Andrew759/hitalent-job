-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employees
(
    id            bigserial PRIMARY KEY,
    department_id BIGINT       NOT NULL,
        CONSTRAINT fk_departments_employees
            FOREIGN KEY (department_id)
            REFERENCES departments(id)
            ON DELETE CASCADE,
    full_name     VARCHAR(200) NOT NULL,
    position      VARCHAR(200) NOT NULL,
    hired_at      TIMESTAMP WITH TIME ZONE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd
