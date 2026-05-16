-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS departments
(
    id         bigserial PRIMARY KEY,
    name       varchar(200) NOT NULL,
    parent_id  bigint,
        CONSTRAINT fk_departments_departments
            FOREIGN KEY (parent_id)
            REFERENCES departments(id)
            ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd
