-- +goose Up
ALTER TABLE expenses
ADD COLUMN amount_usd NUMERIC(12,2);

-- +goose Down
ALTER TABLE expenses
DROP COLUMN amount_usd;
