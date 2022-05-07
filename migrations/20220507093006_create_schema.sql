-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
  log_name VARCHAR(255) PRIMARY KEY,
  head_sn BIGINT
);

CREATE TABLE operations (
  log_name VARCHAR(255) NOT NULL,
  sn BIGINT NOT NULL,
  data BYTEA,
  PRIMARY KEY (log_name, sn),
  CONSTRAINT fk_logs
    FOREIGN KEY (log_name)
    REFERENCES logs(log_name)
    ON DELETE CASCADE
);

ALTER TABLE logs
ADD CONSTRAINT fk_operations
  FOREIGN KEY (log_name, head_sn)
  REFERENCES operations(log_name, sn)
  ON DELETE RESTRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE logs
DROP CONSTRAINT fk_operations;

ALTER TABLE operations
DROP CONSTRAINT fk_logs;

DROP TABLE logs;
DROP TABLE operations;
-- +goose StatementEnd
