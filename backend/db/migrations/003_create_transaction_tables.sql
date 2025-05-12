-- +migrate Up
-- Create data sources table
CREATE TABLE IF NOT EXISTS data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    UNIQUE(name)
);

-- Create match rules table
CREATE TABLE IF NOT EXISTS match_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    match_by_amount BOOLEAN NOT NULL DEFAULT true,
    match_by_date BOOLEAN NOT NULL DEFAULT true,
    date_tolerance INTEGER NOT NULL DEFAULT 0,
    match_by_reference BOOLEAN NOT NULL DEFAULT false,
    active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- Create transaction matches table
CREATE TABLE IF NOT EXISTS transaction_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_status VARCHAR(20) NOT NULL DEFAULT 'Pending' CHECK (match_status IN ('Pending', 'Approved', 'Rejected')),
    match_type VARCHAR(20) NOT NULL CHECK (match_type IN ('Automatic', 'Manual')),
    match_rule_id UUID REFERENCES match_rules(id),
    matched_by UUID NOT NULL REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approval_date TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_source_id UUID NOT NULL REFERENCES data_sources(id),
    transaction_date DATE NOT NULL,
    post_date DATE NOT NULL,
    description TEXT,
    reference VARCHAR(100),
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'Unmatched' CHECK (status IN ('Unmatched', 'Matched', 'Approved')),
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    created_by UUID NOT NULL REFERENCES users(id),
    match_id UUID REFERENCES transaction_matches(id)
);

-- Create transaction upload tracking table
CREATE TABLE IF NOT EXISTS transaction_uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_source_id UUID NOT NULL REFERENCES data_sources(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    upload_date TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    status VARCHAR(20) NOT NULL DEFAULT 'Processing' CHECK (status IN ('Processing', 'Completed', 'Failed')),
    record_count INTEGER DEFAULT 0,
    error_message TEXT
);

-- Add indices for better performance
CREATE INDEX IF NOT EXISTS idx_transactions_data_source_id ON transactions(data_source_id);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_transactions_amount ON transactions(amount);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_match_id ON transactions(match_id);
CREATE INDEX IF NOT EXISTS idx_transaction_matches_status ON transaction_matches(match_status);
CREATE INDEX IF NOT EXISTS idx_match_rules_active ON match_rules(active);
CREATE INDEX IF NOT EXISTS idx_transaction_uploads_status ON transaction_uploads(status);

-- Add audit triggers for transaction changes
CREATE TABLE IF NOT EXISTS transactions_audit (
    id UUID NOT NULL,
    data_source_id UUID NOT NULL,
    transaction_date DATE NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    status VARCHAR(20) NOT NULL,
    match_id UUID,
    operation VARCHAR(10) NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    changed_by UUID
);

CREATE OR REPLACE FUNCTION audit_transaction_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO transactions_audit
        SELECT OLD.id, OLD.data_source_id, OLD.transaction_date, OLD.amount, OLD.status, OLD.match_id, 'DELETE', (NOW() AT TIME ZONE 'UTC'), NULL;
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO transactions_audit
        SELECT NEW.id, NEW.data_source_id, NEW.transaction_date, NEW.amount, NEW.status, NEW.match_id, 'UPDATE', (NOW() AT TIME ZONE 'UTC'), NULL;
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO transactions_audit
        SELECT NEW.id, NEW.data_source_id, NEW.transaction_date, NEW.amount, NEW.status, NEW.match_id, 'INSERT', (NOW() AT TIME ZONE 'UTC'), NULL;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER transactions_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON transactions
FOR EACH ROW EXECUTE FUNCTION audit_transaction_changes();

-- +migrate Down
DROP TRIGGER IF EXISTS transactions_audit_trigger ON transactions;
DROP FUNCTION IF EXISTS audit_transaction_changes();
DROP TABLE IF EXISTS transactions_audit CASCADE;
DROP TABLE IF EXISTS transaction_uploads CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS transaction_matches CASCADE;
DROP TABLE IF EXISTS match_rules CASCADE;
DROP TABLE IF EXISTS data_sources CASCADE; 