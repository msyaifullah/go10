-- Migration Down: Drop loan management schema
-- File: 001_create_loan_schema.down.sql

-- Drop indexes first
DROP INDEX IF EXISTS idx_file_uploads_deleted_at;
DROP INDEX IF EXISTS idx_file_uploads_uploaded_by;
DROP INDEX IF EXISTS idx_file_uploads_entity_type;
DROP INDEX IF EXISTS idx_file_uploads_entity_id;

DROP INDEX IF EXISTS idx_email_notifications_deleted_at;
DROP INDEX IF EXISTS idx_email_notifications_sent_at;
DROP INDEX IF EXISTS idx_email_notifications_loan_id;
DROP INDEX IF EXISTS idx_email_notifications_investor_id;

DROP INDEX IF EXISTS idx_loan_state_histories_deleted_at;
DROP INDEX IF EXISTS idx_loan_state_histories_change_date;
DROP INDEX IF EXISTS idx_loan_state_histories_loan_id;

DROP INDEX IF EXISTS idx_disbursements_deleted_at;
DROP INDEX IF EXISTS idx_disbursements_disbursement_date;
DROP INDEX IF EXISTS idx_disbursements_field_officer_id;
DROP INDEX IF EXISTS idx_disbursements_loan_id;

DROP INDEX IF EXISTS idx_investments_deleted_at;
DROP INDEX IF EXISTS idx_investments_investment_date;
DROP INDEX IF EXISTS idx_investments_investor_id;
DROP INDEX IF EXISTS idx_investments_loan_id;

DROP INDEX IF EXISTS idx_approvals_deleted_at;
DROP INDEX IF EXISTS idx_approvals_approval_date;
DROP INDEX IF EXISTS idx_approvals_validator_id;
DROP INDEX IF EXISTS idx_approvals_loan_id;

DROP INDEX IF EXISTS idx_loans_deleted_at;
DROP INDEX IF EXISTS idx_loans_created_at;
DROP INDEX IF EXISTS idx_loans_state;
DROP INDEX IF EXISTS idx_loans_borrower_id;

DROP INDEX IF EXISTS idx_investors_deleted_at;
DROP INDEX IF EXISTS idx_investors_email;
DROP INDEX IF EXISTS idx_investors_investor_code;

DROP INDEX IF EXISTS idx_employees_deleted_at;
DROP INDEX IF EXISTS idx_employees_role;
DROP INDEX IF EXISTS idx_employees_email;
DROP INDEX IF EXISTS idx_employees_employee_id;

DROP INDEX IF EXISTS idx_borrowers_deleted_at;
DROP INDEX IF EXISTS idx_borrowers_email;
DROP INDEX IF EXISTS idx_borrowers_id_number;

-- Drop tables in correct order (respecting foreign key constraints)
DROP TABLE IF EXISTS file_uploads;
DROP TABLE IF EXISTS email_notifications;
DROP TABLE IF EXISTS loan_state_histories;
DROP TABLE IF EXISTS disbursements;
DROP TABLE IF EXISTS investments;
DROP TABLE IF EXISTS approvals;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS investors;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS borrowers;