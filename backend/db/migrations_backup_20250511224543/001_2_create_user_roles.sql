-- Add user roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('preparer', 'approver', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, role)
);

-- Add index for faster role lookups
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);

-- Create audit table for user roles
CREATE TABLE IF NOT EXISTS user_roles_audit (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL,
    operation VARCHAR(10) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    changed_by UUID
);

-- Trigger for audit trail
CREATE OR REPLACE FUNCTION audit_user_role_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO user_roles_audit
        SELECT OLD.id, OLD.user_id, OLD.role, 'DELETE', NOW(), NULL;
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO user_roles_audit
        SELECT NEW.id, NEW.user_id, NEW.role, 'UPDATE', NOW(), NULL;
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO user_roles_audit
        SELECT NEW.id, NEW.user_id, NEW.role, 'INSERT', NOW(), NULL;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_roles_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON user_roles
FOR EACH ROW EXECUTE FUNCTION audit_user_role_changes(); 