-- Seed data for local development.
-- password_hash below is bcrypt("password") — do NOT use in any real environment.

INSERT INTO users (id, username, email, password_hash, roles) VALUES
    ('00000000-0000-0000-0000-000000000001',
     'admin',
     'admin@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     ARRAY['admin', 'user']),
    ('00000000-0000-0000-0000-000000000002',
     'alice',
     'alice@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     ARRAY['user']),
    ('00000000-0000-0000-0000-000000000003',
     'bob',
     'bob@example.com',
     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
     ARRAY['user'])
ON CONFLICT (id) DO NOTHING;
