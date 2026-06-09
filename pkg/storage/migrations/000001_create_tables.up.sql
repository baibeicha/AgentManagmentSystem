-- 1. Таблица организаций
CREATE TABLE tenants
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица пользователей
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     BIGINT       NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    login         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    global_role   VARCHAR(50)  NOT NULL    DEFAULT 'user',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_global_role CHECK (global_role IN ('global_admin', 'tenant_admin', 'user'))
);

CREATE INDEX idx_users_tenant_id ON users (tenant_id);

-- 3. Группы устройств
CREATE TABLE host_groups
(
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  BIGINT       NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, name)
);

CREATE INDEX idx_host_groups_tenant_id ON host_groups (tenant_id);

-- 4. Устройства (Агенты)
CREATE TABLE hosts
(
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  BIGINT       NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    group_id   BIGINT       REFERENCES host_groups (id) ON DELETE SET NULL,
    hostname   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hosts_tenant_id ON hosts (tenant_id);
CREATE INDEX idx_hosts_group_id ON hosts (group_id);

-- 5. Матрица гранулярных прав (RBAC)
CREATE TABLE group_policies
(
    user_id    BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id   BIGINT       NOT NULL REFERENCES host_groups (id) ON DELETE CASCADE,
    permission VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, group_id, permission)
);

CREATE INDEX idx_group_policies_user_id ON group_policies (user_id);
CREATE INDEX idx_group_policies_group_id ON group_policies (group_id);