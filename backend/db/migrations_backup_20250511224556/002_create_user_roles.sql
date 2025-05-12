-- +migrate Up
-- Add user roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('preparer', 'approver', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
    updated_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
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
    changed_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
    changed_by UUID
);

-- Trigger for audit trail
CREATE OR REPLACE FUNCTION audit_user_role_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO user_roles_audit
        SELECT OLD.id, OLD.user_id, OLD.role, 'DELETE', set_utc_timestamp(), NULL;
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO user_roles_audit
        SELECT NEW.id, NEW.user_id, NEW.role, 'UPDATE', set_utc_timestamp(), NULL;
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO user_roles_audit
        SELECT NEW.id, NEW.user_id, NEW.role, 'INSERT', set_utc_timestamp(), NULL;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_roles_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON user_roles
FOR EACH ROW EXECUTE FUNCTION audit_user_role_changes();

-- +migrate Down
DROP TRIGGER IF EXISTS user_roles_audit_trigger ON user_roles;
DROP FUNCTION IF EXISTS audit_user_role_changes();
DROP TABLE IF EXISTS user_roles_audit CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE; 