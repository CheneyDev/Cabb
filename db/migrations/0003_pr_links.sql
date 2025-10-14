-- PR <-> Issue links
CREATE TABLE IF NOT EXISTS pr_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  plane_issue_id uuid NOT NULL,
  cnb_repo_id text NOT NULL,
  pr_iid text NOT NULL,
  linked_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cnb_repo_id, pr_iid)
);

