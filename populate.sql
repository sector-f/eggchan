INSERT INTO boards (name) VALUES ('diy');
INSERT INTO boards (name) VALUES ('out');

INSERT INTO threads (board_id, comment) VALUES (1, 'first thread on diy');
INSERT INTO threads (board_id, comment) VALUES (2, 'first thread on out');

INSERT INTO comments (reply_to, comment) VALUES (1, 'reply to first thread on diy');
INSERT INTO threads (board_id, comment) VALUES (1, 'new thread on diy');
INSERT INTO comments (reply_to, comment) VALUES (1, 'another reply to first thread on diy');
INSERT INTO comments (reply_to, comment) VALUES (2, 'reply to thread on out');

/* "password" in bcrypt */
INSERT INTO users (username, password) VALUES ('admin', '$2y$10$9bLDFPAILNe5qYmcG1FtmOeCT1dLtUVCU3.rSVzEa782QbPpXSmYy');

INSERT INTO user_permissions (user_id, permission) VALUES
	((SELECT id FROM users WHERE username = 'admin'), (SELECT id FROM permissions WHERE name = 'create_board')),
	((SELECT id FROM users WHERE username = 'admin'), (SELECT id FROM permissions WHERE name = 'delete_board')),
	((SELECT id FROM users WHERE username = 'admin'), (SELECT id FROM permissions WHERE name = 'delete_thread')),
	((SELECT id FROM users WHERE username = 'admin'), (SELECT id FROM permissions WHERE name = 'delete_post'));
