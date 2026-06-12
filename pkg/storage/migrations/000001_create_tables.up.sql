CREATE TABLE tenants
(
    id         UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    name       VARCHAR(255)             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    id             UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    tenant_id      UUID                     NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    email          VARCHAR(255)             NOT NULL UNIQUE,
    password_hash  VARCHAR(255)             NOT NULL,
    role           VARCHAR(50)              NOT NULL DEFAULT 'VIEWER',
    status         VARCHAR(50)              NOT NULL DEFAULT 'ACTIVE',

    is_2fa_enabled BOOLEAN                  NOT NULL DEFAULT FALSE,
    totp_secret    VARCHAR(255),

    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_user_role CHECK (role IN ('TEAM_ADMIN', 'OPERATOR', 'VIEWER')),
    CONSTRAINT chk_user_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'PENDING_INVITE'))
);

CREATE INDEX idx_users_tenant_id ON users (tenant_id);
CREATE INDEX idx_users_email ON users (email);

CREATE TABLE sessions
(
    id            UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    user_id       UUID                     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token VARCHAR(512)             NOT NULL UNIQUE,
    device_id     VARCHAR(255),
    client_ip     VARCHAR(45),
    user_agent    TEXT,
    is_revoked    BOOLEAN                  NOT NULL DEFAULT FALSE,
    expires_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_token ON sessions (refresh_token);

CREATE TABLE resource_policies
(
    id          UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    user_id     UUID                     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    resource_id UUID                     NOT NULL,
    action      VARCHAR(100)             NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (user_id, resource_id, action)
);

CREATE INDEX idx_resource_policies_user_id ON resource_policies (user_id);
CREATE INDEX idx_resource_policies_resource_id ON resource_policies (resource_id);