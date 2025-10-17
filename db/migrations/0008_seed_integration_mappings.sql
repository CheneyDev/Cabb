-- Seed default global mappings to reduce manual setup
-- Priority: Plane priority -> CNB priority

INSERT INTO integration_mappings (
  scope_kind, scope_id, mapping_type,
  left_system, left_type, left_key,
  right_system, right_type, right_key,
  bidirectional, extras, active, created_at, updated_at
)
VALUES
  ('global', NULL, 'priority', 'plane', 'priority', 'urgent', 'cnb', 'priority', 'P0', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'high',   'cnb', 'priority', 'P1', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'medium', 'cnb', 'priority', 'P2', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'low',    'cnb', 'priority', 'P3', true, '{}'::jsonb, true, now(), now()),
  ('global', NULL, 'priority', 'plane', 'priority', 'none',   'cnb', 'priority', ''  , true, '{}'::jsonb, true, now(), now())
ON CONFLICT (scope_kind, scope_id, mapping_type, left_system, left_type, left_key, right_system, right_type, right_key)
DO NOTHING;

