CREATE TABLE IF NOT EXISTS team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(255) NOT NULL, -- gitlab, github, etc
    provider_username VARCHAR(255) NOT NULL,
    slack_user_id VARCHAR(255) NOT NULL,
    seniority_level VARCHAR(255) NOT NULL, -- junior, mid, senior
    weight INTEGER NOT NULL, -- base on the seniority level
    organization_slug VARCHAR(255) NOT NULL, -- provider org: 'nexus-team', etc
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    UNIQUE (provider, provider_username, organization_slug)
);