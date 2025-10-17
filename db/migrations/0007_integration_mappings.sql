-- Unified integration mappings across platforms (Plane, CNB, Lark)
-- Supports multiple mapping types (user/priority/label/state/pr_state/custom)
-- One-to-many/many-to-many supported by inserting multiple rows

CREATE TABLE IF NOT EXISTS integration_mappings (
  id BIGSERIAL PRIMARY KEY,
  scope_kind TEXT NOT NULL CHECK (scope_kind IN ('global','plane_workspace','plane_project','cnb_repo','lark_tenant')),
  scope_id TEXT, -- optional; TEXT to allow various id formats
  mapping_type TEXT NOT NULL CHECK (mapping_type IN ('user','priority','label','state','pr_state','custom')),
  left_system TEXT NOT NULL CHECK (left_system IN ('plane','cnb','lark')),
  left_type TEXT NOT NULL,
  left_key TEXT NOT NULL,
  right_system TEXT NOT NULL CHECK (right_system IN ('plane','cnb','lark')),
  right_type TEXT NOT NULL,
  right_key TEXT NOT NULL,
  bidirectional BOOLEAN NOT NULL DEFAULT true,
  extras JSONB NOT NULL DEFAULT '{}',
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Uniqueness per directional pair in a given scope + mapping_type
DO $$ BEGIN
  ALTER TABLE integration_mappings
    ADD CONSTRAINT uq_integration_map UNIQUE (scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Helpful directional indexes
CREATE INDEX IF NOT EXISTS idx_int_map_left ON integration_mappings(scope_kind, scope_id, mapping_type, left_system, left_type, left_key) WHERE active = true;
CREATE INDEX IF NOT EXISTS idx_int_map_right ON integration_mappings(scope_kind, scope_id, mapping_type, right_system, right_type, right_key) WHERE active = true;

