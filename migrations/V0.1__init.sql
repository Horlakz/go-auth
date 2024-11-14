-- Table for storing roles
CREATE TABLE
    roles (
        id UUID PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT roles_name_unique UNIQUE (name)
    );

-- Table for storing users
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        fullname VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT users_email_unique UNIQUE (email)
    );

-- Table for user roles (many-to-many relation between users and roles)
CREATE TABLE
    user_roles (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        role_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
        CONSTRAINT user_roles_unique UNIQUE (user_id, role_id)
    );

-- Table for verification codes
CREATE TABLE
    verification_codes (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        code VARCHAR(255) NOT NULL,
        purpose VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT verification_codes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
        CONSTRAINT verification_codes_code_unique UNIQUE (user_id, code, purpose)
    );

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);

CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

CREATE INDEX idx_verification_codes_user_id ON verification_codes (user_id);

CREATE INDEX idx_verification_codes_code ON verification_codes (code);