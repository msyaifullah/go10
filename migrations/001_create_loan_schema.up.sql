-- Migration Up: Create loan management schema
-- File: 001_create_loan_schema.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create borrowers table
CREATE TABLE borrowers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_number VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create employees table
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create investors table
CREATE TABLE investors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    investor_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create loans table
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id UUID NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_rate DECIMAL(5,4) NOT NULL,
    roi DECIMAL(5,4) NOT NULL,
    state VARCHAR(20) NOT NULL DEFAULT 'proposed',
    agreement_letter_url TEXT,
    total_invested DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_loans_borrower FOREIGN KEY (borrower_id) REFERENCES borrowers(id),
    CONSTRAINT chk_principal_amount CHECK (principal_amount > 0),
    CONSTRAINT chk_interest_rate CHECK (interest_rate >= 0 AND interest_rate <= 1),
    CONSTRAINT chk_roi CHECK (roi >= 0 AND roi <= 1),
    CONSTRAINT chk_loan_state CHECK (state IN ('proposed', 'approved', 'invested', 'disbursed')),
    CONSTRAINT chk_total_invested CHECK (total_invested >= 0),
    CONSTRAINT chk_total_invested_not_exceed_principal CHECK (total_invested <= principal_amount)
);

-- Create approvals table
CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID UNIQUE NOT NULL,
    validator_id UUID NOT NULL,
    approval_date TIMESTAMP WITH TIME ZONE NOT NULL,
    visit_proof_image_url TEXT NOT NULL,
    visit_proof_image_type VARCHAR(10) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_approvals_loan FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_approvals_validator FOREIGN KEY (validator_id) REFERENCES employees(id),
    CONSTRAINT chk_visit_proof_image_type CHECK (visit_proof_image_type IN ('pdf', 'jpeg', 'png'))
);

-- Create investments table
CREATE TABLE investments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL,
    investor_id UUID NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    investment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expected_return DECIMAL(15,2) NOT NULL,
    agreement_sent BOOLEAN DEFAULT false,
    agreement_sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_investments_loan FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_investments_investor FOREIGN KEY (investor_id) REFERENCES investors(id),
    CONSTRAINT chk_investment_amount CHECK (amount > 0),
    CONSTRAINT chk_expected_return CHECK (expected_return >= 0)
);

-- Create disbursements table
CREATE TABLE disbursements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID UNIQUE NOT NULL,
    field_officer_id UUID NOT NULL,
    disbursement_date TIMESTAMP WITH TIME ZONE NOT NULL,
    signed_agreement_url TEXT NOT NULL,
    signed_agreement_file_type VARCHAR(10) NOT NULL,
    disbursed_amount DECIMAL(15,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_disbursements_loan FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_disbursements_field_officer FOREIGN KEY (field_officer_id) REFERENCES employees(id),
    CONSTRAINT chk_signed_agreement_file_type CHECK (signed_agreement_file_type IN ('pdf', 'jpeg', 'png')),
    CONSTRAINT chk_disbursed_amount CHECK (disbursed_amount > 0)
);

-- Create loan_state_histories table
CREATE TABLE loan_state_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL,
    previous_state VARCHAR(20),
    new_state VARCHAR(20) NOT NULL,
    changed_by UUID NOT NULL,
    change_reason TEXT,
    change_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_loan_state_histories_loan FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT fk_loan_state_histories_changed_by FOREIGN KEY (changed_by) REFERENCES employees(id),
    CONSTRAINT chk_previous_state CHECK (previous_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    CONSTRAINT chk_new_state CHECK (new_state IN ('proposed', 'approved', 'invested', 'disbursed'))
);

-- Create email_notifications table
CREATE TABLE email_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    investor_id UUID NOT NULL,
    loan_id UUID NOT NULL,
    email_type VARCHAR(50) NOT NULL,
    email_subject VARCHAR(255) NOT NULL,
    email_body TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL,
    delivered_at TIMESTAMP WITH TIME ZONE,
    opened_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'sent',
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,


    CONSTRAINT fk_email_notifications_investor FOREIGN KEY (investor_id) REFERENCES investors(id),
    CONSTRAINT fk_email_notifications_loan FOREIGN KEY (loan_id) REFERENCES loans(id),
    CONSTRAINT chk_email_status CHECK (status IN ('sent', 'delivered', 'opened', 'failed'))
);

-- Create file_uploads table
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(10) NOT NULL,
    file_size BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    file_url TEXT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    uploaded_by UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    
    CONSTRAINT fk_file_uploads_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES employees(id),
    CONSTRAINT chk_file_type CHECK (file_type IN ('pdf', 'jpeg', 'png')),
    CONSTRAINT chk_file_size CHECK (file_size > 0)
);

-- Create indexes for better performance
CREATE INDEX idx_borrowers_id_number ON borrowers(id_number);
CREATE INDEX idx_borrowers_email ON borrowers(email);
CREATE INDEX idx_borrowers_deleted_at ON borrowers(deleted_at);

CREATE INDEX idx_employees_employee_id ON employees(employee_id);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_role ON employees(role);
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at);

CREATE INDEX idx_investors_investor_code ON investors(investor_code);
CREATE INDEX idx_investors_email ON investors(email);
CREATE INDEX idx_investors_deleted_at ON investors(deleted_at);

CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_state ON loans(state);
CREATE INDEX idx_loans_created_at ON loans(created_at);
CREATE INDEX idx_loans_deleted_at ON loans(deleted_at);

CREATE INDEX idx_approvals_loan_id ON approvals(loan_id);
CREATE INDEX idx_approvals_validator_id ON approvals(validator_id);
CREATE INDEX idx_approvals_approval_date ON approvals(approval_date);
CREATE INDEX idx_approvals_deleted_at ON approvals(deleted_at);

CREATE INDEX idx_investments_loan_id ON investments(loan_id);
CREATE INDEX idx_investments_investor_id ON investments(investor_id);
CREATE INDEX idx_investments_investment_date ON investments(investment_date);
CREATE INDEX idx_investments_deleted_at ON investments(deleted_at);

CREATE INDEX idx_disbursements_loan_id ON disbursements(loan_id);
CREATE INDEX idx_disbursements_field_officer_id ON disbursements(field_officer_id);
CREATE INDEX idx_disbursements_disbursement_date ON disbursements(disbursement_date);
CREATE INDEX idx_disbursements_deleted_at ON disbursements(deleted_at);

CREATE INDEX idx_loan_state_histories_loan_id ON loan_state_histories(loan_id);
CREATE INDEX idx_loan_state_histories_change_date ON loan_state_histories(change_date);
CREATE INDEX idx_loan_state_histories_deleted_at ON loan_state_histories(deleted_at);

CREATE INDEX idx_email_notifications_investor_id ON email_notifications(investor_id);
CREATE INDEX idx_email_notifications_loan_id ON email_notifications(loan_id);
CREATE INDEX idx_email_notifications_sent_at ON email_notifications(sent_at);
CREATE INDEX idx_email_notifications_deleted_at ON email_notifications(deleted_at);

CREATE INDEX idx_file_uploads_entity_id ON file_uploads(entity_id);
CREATE INDEX idx_file_uploads_entity_type ON file_uploads(entity_type);
CREATE INDEX idx_file_uploads_uploaded_by ON file_uploads(uploaded_by);
CREATE INDEX idx_file_uploads_deleted_at ON file_uploads(deleted_at);