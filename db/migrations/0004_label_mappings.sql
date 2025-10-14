-- Label mappings between CNB and Plane
CREATE TABLE IF NOT EXISTS label_mappings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_project_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  cnb_label text NOT NULL,
  plane_label_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (plane_project_id, cnb_repo_id, cnb_label)
);

