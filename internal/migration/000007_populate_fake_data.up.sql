INSERT INTO  users (id, name, last_name, username, email, password, created_at, updated_at)
VALUES
	(1, 'John', 'Doe', 'john_doe', 'john@example.com', crypt('password123',gen_salt('bf')), NOW(), NOW()),
	(2, 'Jane', 'Smith', 'jane_smith', 'jane@example.com', crypt('password123',gen_salt('bf')), NOW(), NOW()),
	(3, 'Alice', 'Johnson', 'alice_johnson', 'alice@example.com', crypt('password123',gen_salt('bf')), NOW(), NOW());

INSERT INTO follows (user_id, follow_id)
VALUES
	(1, 2),
	(1, 3),
	(2, 1),
	(2, 3),
	(3, 1),
	(3, 2);

INSERT INTO posts (user_id, content)
VALUES
	(1, 'Hello World'),
	(2, 'Hello Go'),
	(3, 'Hello Gin');


INSERT INTO comments (id, post_id, user_id, content)
VALUES
	(gen_random_uuid(),1, 1, 'Great post!'),
	(gen_random_uuid(),1, 2, 'Thanks for sharing!'),
	(gen_random_uuid(),1, 3, 'Interesting perspective.'),
	(gen_random_uuid(),2, 3, 'I learned something new.'),
	(gen_random_uuid(),2, 2, 'Well written!'),
	(gen_random_uuid(),3, 3, 'I completely agree.');


