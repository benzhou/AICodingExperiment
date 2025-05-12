-- Create match sets table
CREATE TABLE IF NOT EXISTS match_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES rules(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name, tenant_id)
);

-- Create match set data sources table (junction table for many-to-many relationship)
CREATE TABLE IF NOT EXISTS match_set_data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(match_set_id, data_source_id)
);

-- Create matched transactions table
CREATE TABLE IF NOT EXISTS matched_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    match_group_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(transaction_id, match_set_id)
);

-- Create unmatched transactions table
CREATE TABLE IF NOT EXISTS unmatched_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(transaction_id, match_set_id)
);

-- Create match progress table
CREATE TABLE IF NOT EXISTS match_progress (
    match_set_id UUID PRIMARY KEY REFERENCES match_sets(id) ON DELETE CASCADE,
    total_transactions INTEGER NOT NULL DEFAULT 0,
    processed_transactions INTEGER NOT NULL DEFAULT 0,
    matched_transactions INTEGER NOT NULL DEFAULT 0,
    unmatched_transactions INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'Running', 'Completed', 'Failed')),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error TEXT
);

-- Create indices for better performance
CREATE INDEX IF NOT EXISTS idx_match_sets_tenant_id ON match_sets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_match_sets_rule_id ON match_sets(rule_id);
CREATE INDEX IF NOT EXISTS idx_match_set_data_sources_match_set_id ON match_set_data_sources(match_set_id);
CREATE INDEX IF NOT EXISTS idx_match_set_data_sources_data_source_id ON match_set_data_sources(data_source_id);
CREATE INDEX IF NOT EXISTS idx_matched_transactions_match_set_id ON matched_transactions(match_set_id);
CREATE INDEX IF NOT EXISTS idx_matched_transactions_match_group_id ON matched_transactions(match_group_id);
CREATE INDEX IF NOT EXISTS idx_matched_transactions_tenant_id ON matched_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_unmatched_transactions_match_set_id ON unmatched_transactions(match_set_id);
CREATE INDEX IF NOT EXISTS idx_unmatched_transactions_tenant_id ON unmatched_transactions(tenant_id); 