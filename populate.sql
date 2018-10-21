INSERT INTO boards (name) VALUES ('diy');
INSERT INTO boards (name) VALUES ('out');

INSERT INTO threads (board_id, comment) VALUES (1, 'first thread on diy');
INSERT INTO threads (board_id, comment) VALUES (2, 'first thread on out');

INSERT INTO posts (reply_to, comment) VALUES (1, 'reply to first thread on diy');
INSERT INTO threads (board_id, comment) VALUES (1, 'new thread on diy');
INSERT INTO posts (reply_to, comment) VALUES (1, 'another reply to first thread on diy');
INSERT INTO posts (reply_to, comment) VALUES (2, 'reply to thread on out');
