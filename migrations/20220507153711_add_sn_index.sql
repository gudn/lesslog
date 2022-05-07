-- +goose Up
-- +goose StatementBegin
CREATE INDEX log_sn_idx ON operations(log_name, sn ASC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX log_sn_idx;
-- +goose StatementEnd
