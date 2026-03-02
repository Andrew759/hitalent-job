-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS departments
(
    id         bigserial PRIMARY KEY,
    name       varchar(200) NOT NULL,
    parent_id  bigint,
        CONSTRAINT fk_departments_departments
            REFERENCES departments
            ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE
);


-- TODO: подумать, возможно стоит его использовать
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_name_parent
--     ON departments (name, (COALESCE(parent_id, -1)));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd
