-- Users table
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- OAuth2 Clients
CREATE TABLE IF NOT EXISTS oauth_clients (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id       TEXT UNIQUE NOT NULL,
    client_secret   TEXT NOT NULL,
    redirect_uris   TEXT[] NOT NULL,
    scopes          TEXT,
    grant_types     TEXT[],
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- OAuth2 Access & Refresh Tokens
CREATE TABLE IF NOT EXISTS oauth_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    client_id       UUID REFERENCES oauth_clients(id) ON DELETE CASCADE,
    access_token    TEXT UNIQUE NOT NULL,
    refresh_token   TEXT UNIQUE,
    scopes          TEXT,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Authorization Codes
CREATE TABLE IF NOT EXISTS oauth_authorization_codes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            TEXT UNIQUE NOT NULL,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    client_id       UUID REFERENCES oauth_clients(id) ON DELETE CASCADE,
    redirect_uri    TEXT NOT NULL,
    scopes          TEXT,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Optional scopes registry
CREATE TABLE IF NOT EXISTS oauth_scopes (
    name        VARCHAR(50) PRIMARY KEY,
    description TEXT
);
