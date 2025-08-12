-- Seed Data: Sample data to simulate complete loan cycles
-- File: 002_seed_data.up.sql

-- Insert sample borrowers
INSERT INTO borrowers (id, id_number, first_name, last_name, email, phone_number, address) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'ID001234567', 'John', 'Doe', 'john.doe@email.com', '+62812345678', 'Jl. Sudirman No. 123, Jakarta'),
('550e8400-e29b-41d4-a716-446655440002', 'ID001234568', 'Jane', 'Smith', 'jane.smith@email.com', '+62812345679', 'Jl. Thamrin No. 456, Jakarta'),
('550e8400-e29b-41d4-a716-446655440003', 'ID001234569', 'Bob', 'Johnson', 'bob.johnson@email.com', '+62812345680', 'Jl. Gatot Subroto No. 789, Jakarta'),
('550e8400-e29b-41d4-a716-446655440004', 'ID001234570', 'Alice', 'Brown', 'alice.brown@email.com', '+62812345681', 'Jl. Rasuna Said No. 321, Jakarta'),
('550e8400-e29b-41d4-a716-446655440005', 'ID001234571', 'Charlie', 'Wilson', 'charlie.wilson@email.com', '+62812345682', 'Jl. HR Rasuna Said No. 654, Jakarta');

-- Insert sample employees
INSERT INTO employees (id, employee_id, first_name, last_name, email, role, phone_number, is_active) VALUES
('a0000000-0000-0000-0000-000000000001', 'SYS001', 'System', 'Employee', 'system@company.com', 'system', '+62800000000', true),
('660e8400-e29b-41d4-a716-446655440001', 'EMP001', 'Sarah', 'Anderson', 'sarah.anderson@company.com', 'field_validator', '+62813456789', true),
('660e8400-e29b-41d4-a716-446655440002', 'EMP002', 'Mike', 'Taylor', 'mike.taylor@company.com', 'field_validator', '+62813456790', true),
('660e8400-e29b-41d4-a716-446655440003', 'EMP003', 'Lisa', 'Davis', 'lisa.davis@company.com', 'field_officer', '+62813456791', true),
('660e8400-e29b-41d4-a716-446655440004', 'EMP004', 'Tom', 'Miller', 'tom.miller@company.com', 'field_officer', '+62813456792', true);

-- Insert sample investors
INSERT INTO investors (id, investor_code, name, email, phone_number, is_active) VALUES
('770e8400-e29b-41d4-a716-446655440001', 'INV001', 'Global Investment Fund', 'contact@globalinvest.com', '+62814567890', true),
('770e8400-e29b-41d4-a716-446655440002', 'INV002', 'Jakarta Capital Partners', 'info@jakartacapital.com', '+62814567891', true),
('770e8400-e29b-41d4-a716-446655440003', 'INV003', 'Prosperity Investment Group', 'hello@prosperitygroup.com', '+62814567892', true),
('770e8400-e29b-41d4-a716-446655440004', 'INV004', 'Indonesia Growth Fund', 'contact@indonesiagrowth.com', '+62814567893', true),
('770e8400-e29b-41d4-a716-446655440005', 'INV005', 'Asia Pacific Ventures', 'info@asiapacificventures.com', '+62814567894', true);

-- Insert sample loans (covering complete lifecycle and different states)
INSERT INTO loans (id, borrower_id, principal_amount, interest_rate, roi, state, agreement_letter_url, total_invested, created_at) VALUES
-- Loan 1: Complete cycle (disbursed)
('880e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 100000000.00, 0.12, 0.10, 'disbursed', 'https://storage.example.com/agreements/loan-001.pdf', 100000000.00, '2024-01-15 10:00:00+07'),

-- Loan 2: Complete cycle (disbursed)
('880e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 75000000.00, 0.15, 0.12, 'disbursed', 'https://storage.example.com/agreements/loan-002.pdf', 75000000.00, '2024-02-01 09:30:00+07'),

-- Loan 3: In invested state (ready for disbursement)
('880e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 50000000.00, 0.18, 0.15, 'invested', 'https://storage.example.com/agreements/loan-003.pdf', 50000000.00, '2024-03-10 14:15:00+07'),

-- Loan 4: In approved state (available for investment)
('880e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440004', 80000000.00, 0.14, 0.11, 'approved', 'https://storage.example.com/agreements/loan-004.pdf', 30000000.00, '2024-03-20 11:45:00+07'),

-- Loan 5: In proposed state (pending approval)
('880e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440005', 120000000.00, 0.16, 0.13, 'proposed', 'https://storage.example.com/agreements/loan-005.pdf', 0.00, '2024-03-25 16:20:00+07');

-- Insert approvals for approved, invested, and disbursed loans
INSERT INTO approvals (id, loan_id, validator_id, approval_date, visit_proof_image_url, visit_proof_image_type, notes) VALUES
('990e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', '2024-01-16 14:30:00+07', 'https://storage.example.com/proof/visit-001.jpeg', 'jpeg', 'Borrower location verified, business operational'),
('990e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440002', '2024-02-02 10:15:00+07', 'https://storage.example.com/proof/visit-002.jpeg', 'jpeg', 'Small business validated, good cash flow'),
('990e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440001', '2024-03-11 09:45:00+07', 'https://storage.example.com/proof/visit-003.png', 'png', 'Manufacturing business, equipment verified'),
('990e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440002', '2024-03-21 13:20:00+07', 'https://storage.example.com/proof/visit-004.jpeg', 'jpeg', 'Retail business, good location and inventory');

-- Insert investments for invested and disbursed loans
INSERT INTO investments (id, loan_id, investor_id, amount, investment_date, expected_return, agreement_sent, agreement_sent_at) VALUES
-- Investments for Loan 1 (disbursed)
('aa0e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440001', 40000000.00, '2024-01-17 10:00:00+07', 4000000.00, true, '2024-01-18 15:30:00+07'),
('aa0e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440002', 35000000.00, '2024-01-18 11:30:00+07', 3500000.00, true, '2024-01-18 15:30:00+07'),
('aa0e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440003', 25000000.00, '2024-01-19 14:15:00+07', 2500000.00, true, '2024-01-19 16:45:00+07'),

-- Investments for Loan 2 (disbursed)
('aa0e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440002', 45000000.00, '2024-02-03 09:20:00+07', 5400000.00, true, '2024-02-04 14:00:00+07'),
('aa0e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440004', 30000000.00, '2024-02-04 13:45:00+07', 3600000.00, true, '2024-02-04 14:00:00+07'),

-- Investments for Loan 3 (invested)
('aa0e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440001', 20000000.00, '2024-03-12 10:30:00+07', 3000000.00, true, '2024-03-13 09:15:00+07'),
('aa0e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440003', 30000000.00, '2024-03-13 15:20:00+07', 4500000.00, true, '2024-03-13 09:15:00+07'),

-- Investments for Loan 4 (approved - partially invested)
('aa0e8400-e29b-41d4-a716-446655440008', '880e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440005', 30000000.00, '2024-03-22 11:00:00+07', 3300000.00, false, NULL);

-- Insert disbursements for disbursed loans
INSERT INTO disbursements (id, loan_id, field_officer_id, disbursement_date, signed_agreement_url, signed_agreement_file_type, disbursed_amount, notes) VALUES
('bb0e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440003', '2024-01-20 10:30:00+07', 'https://storage.example.com/signed/agreement-001.pdf', 'pdf', 100000000.00, 'Funds disbursed successfully to borrower bank account'),
('bb0e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440004', '2024-02-05 14:45:00+07', 'https://storage.example.com/signed/agreement-002.pdf', 'pdf', 75000000.00, 'Disbursement completed, borrower confirmed receipt');

-- Insert loan state history for all state transitions
INSERT INTO loan_state_histories (id, loan_id, previous_state, new_state, changed_by, change_reason, change_date) VALUES
-- Loan 1 history (complete cycle)
('cc0e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 'proposed', 'approved', '660e8400-e29b-41d4-a716-446655440001', 'Field validation completed successfully', '2024-01-16 14:30:00+07'),
('cc0e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440001', 'approved', 'invested', 'a0000000-0000-0000-0000-000000000001', 'Full funding reached from investors', '2024-01-19 14:15:00+07'),
('cc0e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440001', 'invested', 'disbursed', '660e8400-e29b-41d4-a716-446655440003', 'Loan disbursed to borrower', '2024-01-20 10:30:00+07'),

-- Loan 2 history (complete cycle)
('cc0e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440002', 'proposed', 'approved', '660e8400-e29b-41d4-a716-446655440002', 'Business validation successful', '2024-02-02 10:15:00+07'),
('cc0e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440002', 'approved', 'invested', 'a0000000-0000-0000-0000-000000000001', 'Investment target achieved', '2024-02-04 13:45:00+07'),
('cc0e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440002', 'invested', 'disbursed', '660e8400-e29b-41d4-a716-446655440004', 'Funds successfully disbursed', '2024-02-05 14:45:00+07'),

-- Loan 3 history (to invested state)
('cc0e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440003', 'proposed', 'approved', '660e8400-e29b-41d4-a716-446655440001', 'Manufacturing business approved after site visit', '2024-03-11 09:45:00+07'),
('cc0e8400-e29b-41d4-a716-446655440008', '880e8400-e29b-41d4-a716-446655440003', 'approved', 'invested', 'a0000000-0000-0000-0000-000000000001', 'Full investment amount secured', '2024-03-13 15:20:00+07'),

-- Loan 4 history (to approved state)
('cc0e8400-e29b-41d4-a716-446655440009', '880e8400-e29b-41d4-a716-446655440004', 'proposed', 'approved', '660e8400-e29b-41d4-a716-446655440002', 'Retail business approved, awaiting full investment', '2024-03-21 13:20:00+07');

-- Insert email notifications for invested and disbursed loans
INSERT INTO email_notifications (id, investor_id, loan_id, email_type, email_subject, email_body, sent_at, delivered_at, opened_at, status) VALUES
-- Notifications for Loan 1 investors
('dd0e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 'agreement_notification', 'Loan Investment Agreement - Loan #001', 'Dear Investor, Your investment in Loan #001 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-001.pdf', '2024-01-18 15:30:00+07', '2024-01-18 15:35:00+07', '2024-01-18 16:20:00+07', 'opened'),
('dd0e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440001', 'agreement_notification', 'Loan Investment Agreement - Loan #001', 'Dear Investor, Your investment in Loan #001 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-001.pdf', '2024-01-18 15:30:00+07', '2024-01-18 15:37:00+07', '2024-01-18 17:45:00+07', 'opened'),
('dd0e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440001', 'agreement_notification', 'Loan Investment Agreement - Loan #001', 'Dear Investor, Your investment in Loan #001 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-001.pdf', '2024-01-19 16:45:00+07', '2024-01-19 16:50:00+07', '2024-01-19 18:30:00+07', 'opened'),

-- Notifications for Loan 2 investors
('dd0e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440002', 'agreement_notification', 'Loan Investment Agreement - Loan #002', 'Dear Investor, Your investment in Loan #002 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-002.pdf', '2024-02-04 14:00:00+07', '2024-02-04 14:05:00+07', '2024-02-04 15:30:00+07', 'opened'),
('dd0e8400-e29b-41d4-a716-446655440005', '770e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440002', 'agreement_notification', 'Loan Investment Agreement - Loan #002', 'Dear Investor, Your investment in Loan #002 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-002.pdf', '2024-02-04 14:00:00+07', '2024-02-04 14:03:00+07', NULL, 'delivered'),

-- Notifications for Loan 3 investors
('dd0e8400-e29b-41d4-a716-446655440006', '770e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440003', 'agreement_notification', 'Loan Investment Agreement - Loan #003', 'Dear Investor, Your investment in Loan #003 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-003.pdf', '2024-03-13 09:15:00+07', '2024-03-13 09:20:00+07', '2024-03-13 10:45:00+07', 'opened'),
('dd0e8400-e29b-41d4-a716-446655440007', '770e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440003', 'agreement_notification', 'Loan Investment Agreement - Loan #003', 'Dear Investor, Your investment in Loan #003 has been fully funded. Please find the agreement letter attached. Click here to access: https://storage.example.com/agreements/loan-003.pdf', '2024-03-13 09:15:00+07', '2024-03-13 09:18:00+07', '2024-03-13 14:20:00+07', 'opened');

-- Insert sample file uploads for documentation
INSERT INTO file_uploads (id, file_name, file_type, file_size, file_path, file_url, content_type, uploaded_by, entity_type, entity_id, is_active) VALUES
-- Visit proof images for approvals
('ee0e8400-e29b-41d4-a716-446655440001', 'visit_proof_loan_001.jpeg', 'jpeg', 2048576, '/uploads/proof/visit_proof_loan_001.jpeg', 'https://storage.example.com/proof/visit-001.jpeg', 'image/jpeg', '660e8400-e29b-41d4-a716-446655440001', 'approval', '990e8400-e29b-41d4-a716-446655440001', true),
('ee0e8400-e29b-41d4-a716-446655440002', 'visit_proof_loan_002.jpeg', 'jpeg', 1875456, '/uploads/proof/visit_proof_loan_002.jpeg', 'https://storage.example.com/proof/visit-002.jpeg', 'image/jpeg', '660e8400-e29b-41d4-a716-446655440002', 'approval', '990e8400-e29b-41d4-a716-446655440002', true),
('ee0e8400-e29b-41d4-a716-446655440003', 'visit_proof_loan_003.png', 'png', 3145728, '/uploads/proof/visit_proof_loan_003.png', 'https://storage.example.com/proof/visit-003.png', 'image/png', '660e8400-e29b-41d4-a716-446655440001', 'approval', '990e8400-e29b-41d4-a716-446655440003', true),
('ee0e8400-e29b-41d4-a716-446655440004', 'visit_proof_loan_004.jpeg', 'jpeg', 2234567, '/uploads/proof/visit_proof_loan_004.jpeg', 'https://storage.example.com/proof/visit-004.jpeg', 'image/jpeg', '660e8400-e29b-41d4-a716-446655440002', 'approval', '990e8400-e29b-41d4-a716-446655440004', true),

-- Signed agreement documents for disbursements
('ee0e8400-e29b-41d4-a716-446655440005', 'signed_agreement_loan_001.pdf', 'pdf', 5242880, '/uploads/agreements/signed_agreement_loan_001.pdf', 'https://storage.example.com/signed/agreement-001.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440003', 'disbursement', 'bb0e8400-e29b-41d4-a716-446655440001', true),
('ee0e8400-e29b-41d4-a716-446655440006', 'signed_agreement_loan_002.pdf', 'pdf', 4718592, '/uploads/agreements/signed_agreement_loan_002.pdf', 'https://storage.example.com/signed/agreement-002.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440004', 'disbursement', 'bb0e8400-e29b-41d4-a716-446655440002', true),

-- Agreement letter templates for loans
('ee0e8400-e29b-41d4-a716-446655440007', 'agreement_letter_loan_001.pdf', 'pdf', 3932160, '/uploads/agreements/agreement_letter_loan_001.pdf', 'https://storage.example.com/agreements/loan-001.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440001', 'loan', '880e8400-e29b-41d4-a716-446655440001', true),
('ee0e8400-e29b-41d4-a716-446655440008', 'agreement_letter_loan_002.pdf', 'pdf', 4194304, '/uploads/agreements/agreement_letter_loan_002.pdf', 'https://storage.example.com/agreements/loan-002.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440002', 'loan', '880e8400-e29b-41d4-a716-446655440002', true),
('ee0e8400-e29b-41d4-a716-446655440009', 'agreement_letter_loan_003.pdf', 'pdf', 3670016, '/uploads/agreements/agreement_letter_loan_003.pdf', 'https://storage.example.com/agreements/loan-003.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440001', 'loan', '880e8400-e29b-41d4-a716-446655440003', true),
('ee0e8400-e29b-41d4-a716-446655440010', 'agreement_letter_loan_004.pdf', 'pdf', 4456448, '/uploads/agreements/agreement_letter_loan_004.pdf', 'https://storage.example.com/agreements/loan-004.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440002', 'loan', '880e8400-e29b-41d4-a716-446655440004', true),
('ee0e8400-e29b-41d4-a716-446655440011', 'agreement_letter_loan_005.pdf', 'pdf', 3801088, '/uploads/agreements/agreement_letter_loan_005.pdf', 'https://storage.example.com/agreements/loan-005.pdf', 'application/pdf', '660e8400-e29b-41d4-a716-446655440001', 'loan', '880e8400-e29b-41d4-a716-446655440005', true);

-- Update timestamps for more recent data
UPDATE loans SET updated_at = NOW() WHERE state IN ('approved', 'invested');
UPDATE investments SET created_at = investment_date, updated_at = investment_date;
UPDATE approvals SET created_at = approval_date, updated_at = approval_date;
UPDATE disbursements SET created_at = disbursement_date, updated_at = disbursement_date;
UPDATE loan_state_histories SET created_at = change_date, updated_at = change_date;
UPDATE email_notifications SET created_at = sent_at, updated_at = COALESCE(opened_at, delivered_at, sent_at);