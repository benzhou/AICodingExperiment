-- +migrate Up
-- Create match progress table
CREATE TABLE IF NOT EXISTS match_progress (
    match_set_id UUID PRIMARY KEY REFERENCES match_sets(id) ON DELETE CASCADE,
    total_transactions INTEGER NOT NULL DEFAULT 0,
    processed_transactions INTEGER NOT NULL DEFAULT 0,
    matched_transactions INTEGER NOT NULL DEFAULT 0,
    unmatched_transactions INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'Running', 'Completed', 'Failed')),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error TEXT
);

-- Create indices for match_sets table if not already created
CREATE INDEX IF NOT EXISTS idx_match_sets_rule_id ON match_sets(rule_id);
CREATE INDEX IF NOT EXISTS idx_match_set_data_sources_match_set_id ON match_set_data_sources(match_set_id);
CREATE INDEX IF NOT EXISTS idx_match_set_data_sources_data_source_id ON match_set_data_sources(data_source_id);

-- +migrate Down
DROP TABLE IF EXISTS match_progress CASCADE; 