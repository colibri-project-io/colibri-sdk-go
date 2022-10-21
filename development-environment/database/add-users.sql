INSERT INTO profiles (id, name) VALUES (100, 'ADMIN');
INSERT INTO profiles (id, name) VALUES (200, 'USER');

INSERT INTO users (id, name, birthday, profile_id) VALUES (1, 'ADMIN USER', '2022-01-10', 100);
INSERT INTO users (id, name, birthday, profile_id) VALUES (2, 'OTHER USER', '2020-03-25', 200);
