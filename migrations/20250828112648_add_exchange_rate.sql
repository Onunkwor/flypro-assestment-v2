-- +goose Up

ALTER TABLE expenses
ADD COLUMN exchange_rate FLOAT;


-- +goose Down
ALTER TABLE expenses
DROP COLUMN exchange_rate;
