-- Seed admin user  login: admin@dekanat.local  password: admin123
INSERT INTO users (username, email, password_hash, status, is_active)
VALUES ('admin', 'admin@dekanat.local', '$2a$12$kqEPuG/Uf9w598F9RFuQdOUpmixMQf.dY7.J/rlXfaN7BfO4nRa.e', 'active', true)
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.email = 'admin@dekanat.local' AND r.name = 'admin'
ON CONFLICT DO NOTHING;
