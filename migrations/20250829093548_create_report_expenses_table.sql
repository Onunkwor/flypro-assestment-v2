-- +goose Up
CREATE TABLE report_expenses (
    report_id INT NOT NULL REFERENCES expense_reports(id) ON DELETE CASCADE,
    expense_id INT NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    PRIMARY KEY (report_id, expense_id)
);

-- +goose Down
DROP TABLE report_expenses;
