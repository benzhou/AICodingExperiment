-- Create tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(name)
);

-- Create tenant_users table to associate users with tenants
CREATE TABLE IF NOT EXISTS tenant_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, user_id)
);

-- Create role_permissions table for action-based permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) NOT NULL CHECK (role_name IN ('preparer', 'approver', 'admin')),
    permission VARCHAR(50) NOT NULL,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create a unique constraint on role, permission, and tenant
CREATE UNIQUE INDEX idx_role_perm_tenant ON role_permissions(role_name, permission, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'));

-- Create match_sets table
CREATE TABLE IF NOT EXISTS match_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES match_rules(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create match_set_data_sources table to associate data sources with match sets
CREATE TABLE IF NOT EXISTS match_set_data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(match_set_id, data_source_id)
);

-- Add tenant_id to data_sources table
ALTER TABLE data_sources ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- Add tenant_id to match_rules table
ALTER TABLE match_rules ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- Add tenant_id to transactions table
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- Add indices for better performance
CREATE INDEX IF NOT EXISTS idx_data_sources_tenant_id ON data_sources(tenant_id);
CREATE INDEX IF NOT EXISTS idx_match_rules_tenant_id ON match_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tenant_id ON transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_match_sets_tenant_id ON match_sets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_users_tenant_id ON tenant_users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_users_user_id ON tenant_users(user_id);

-- Create tables for matched and unmatched transactions
CREATE TABLE IF NOT EXISTS matched_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    match_group_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(transaction_id)
);

CREATE TABLE IF NOT EXISTS unmatched_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_set_id UUID NOT NULL REFERENCES match_sets(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    reason TEXT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(transaction_id, match_set_id)
);

-- Create indices for transaction tables
CREATE INDEX IF NOT EXISTS idx_matched_transactions_tenant_id ON matched_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_matched_transactions_match_set_id ON matched_transactions(match_set_id);
CREATE INDEX IF NOT EXISTS idx_matched_transactions_match_group_id ON matched_transactions(match_group_id);
CREATE INDEX IF NOT EXISTS idx_unmatched_transactions_tenant_id ON unmatched_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_unmatched_transactions_match_set_id ON unmatched_transactions(match_set_id); 