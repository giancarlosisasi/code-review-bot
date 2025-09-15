CREATE TABLE IF NOT EXISTS mr_assigments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider VARCHAR(255) NOT NULL, -- gitlab, github, etc
    organization_slug VARCHAR(255) NOT NULL, -- provider org: 'nexus-team', etc
    project_path VARCHAR(255) NOT NULL, -- provider project: 'nexus-team/nexus-team-app', etc
    mr_id VARCHAR(255) NOT NULL, -- provider mr: '1234567890', etc
    mr_url VARCHAR(255) NOT NULL, -- provider mr url: 'https://gitlab.com/nexus-team/nexus-team-app/merge_requests/1234567890', etc
    title VARCHAR(255) NOT NULL, -- provider mr title: 'Add new feature', etc
    author_username VARCHAR(255) NOT NULL, -- provider mr author username: 'john_doe', etc
    assignee_id INTEGER REFERENCES team_members(id),
    changes_count INTEGER DEFAULT 0,
    status VARCHAR(255) NOT NULL, -- open, closed, merged, etc
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    UNIQUE (provider, organization_slug, project_path, mr_id, assignee_id)
);

CREATE INDEX idx_mr_assignments_status ON mr_assignments(status);
CREATE INDEX idx_mr_assignments_assignee ON mr_assignments(assignee_id);