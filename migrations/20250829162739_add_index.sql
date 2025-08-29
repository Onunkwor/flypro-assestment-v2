-- +goose Up

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_expense_reports_user_id ON expense_reports (user_id);
CREATE INDEX IF NOT EXISTS idx_expense_reports_status ON expense_reports (status);

CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses (user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses (category);
CREATE INDEX IF NOT EXISTS idx_expenses_id_user_id ON expenses (id, user_id);

CREATE INDEX IF NOT EXISTS idx_report_expenses_report_id ON report_expenses (report_id);
CREATE INDEX IF NOT EXISTS idx_report_expenses_expense_id ON report_expenses (expense_id);

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_expense_reports_user_id;
DROP INDEX IF EXISTS idx_expense_reports_status;

DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_expenses_category;
DROP INDEX IF EXISTS idx_expenses_id_user_id;

DROP INDEX IF EXISTS idx_report_expenses_report_id;
DROP INDEX IF EXISTS idx_report_expenses_expense_id;
-- +goose StatementEnd
