BEGIN;

CREATE TABLE IF NOT EXISTS categories (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS boards (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT NOT NULL,
	category INTEGER NOT NULL REFERENCES categories,
	bump_limit INTEGER NOT NULL DEFAULT 300,
	post_limit INTEGER NOT NULL DEFAULT 500,
	max_num_threads INTEGER NOT NULL DEFAULT 30
);

/* Table board_postnum keeps track of the highest post number on each board */
CREATE TABLE IF NOT EXISTS board_postnum (
	board_id INTEGER PRIMARY KEY REFERENCES boards,
	postnum INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads (
	id SERIAL PRIMARY KEY,
	subject TEXT,
	author TEXT DEFAULT 'Anonymous',
	post_num INTEGER,
	board_id INTEGER REFERENCES boards,
	time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	comment TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
	id SERIAL PRIMARY KEY,
	author TEXT DEFAULT 'Anonymous',
	post_num INTEGER,
	reply_to INTEGER REFERENCES threads ON DELETE CASCADE,
	time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	comment TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

INSERT INTO permissions (name) VALUES
	('create_board'),
	('delete_board'),
	('delete_thread'),
	('delete_post');

CREATE TABLE IF NOT EXISTS user_permissions (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users ON DELETE CASCADE,
	permission INT NOT NULL REFERENCES permissions ON DELETE CASCADE,
	UNIQUE (user_id, permission)
);

CREATE VIEW last_reply_time AS
SELECT ROW_NUMBER() OVER(PARTITION BY boards.id ORDER BY
			CASE
				WHEN MAX(comments.time) IS NOT NULL AND COUNT(*) >= (SELECT bump_limit FROM boards WHERE boards.id = threads.board_id) THEN (SELECT comments.time FROM comments OFFSET (SELECT bump_limit FROM boards WHERE boards.id = threads.board_id) LIMIT 1)
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END DESC) AS sort_number, boards.id AS board_id, boards.max_num_threads, threads.id AS thread_id
FROM threads
LEFT JOIN comments ON comments.reply_to = threads.id
LEFT JOIN boards ON boards.id = threads.board_id
GROUP BY boards.id, threads.id
ORDER BY boards.id, sort_number;

/* CREATE VIEW comments_with_board_info AS */
/* SELECT */
/* 	b.id AS board_id, b.name AS board_name, b.max_num_threads AS board_max_num_threads, */
/* 	c.id, c.author, c.post_num, c.reply_to, c.image, c.time, c.comment */
/* FROM boards b */
/* INNER JOIN threads t ON t.board_id = b.id */
/* INNER JOIN comments c ON c.reply_to = t.id */
/* GROUP BY b.id, c.id */
/* ORDER BY b.id, c.post_num; */

/*****************/

CREATE FUNCTION thread_lock_check_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	BEGIN
		IF
			(SELECT COUNT(*) FROM comments WHERE reply_to = NEW.reply_to) >= (SELECT b.post_limit FROM boards b WHERE b.id = (SELECT b.id FROM boards b INNER JOIN threads t ON b.id = t.board_id WHERE t.id = NEW.reply_to))
		THEN
			RAISE EXCEPTION 'Thread has reached post limit';
		END IF;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER thread_lock_check
	BEFORE INSERT ON comments
	FOR EACH ROW
	EXECUTE PROCEDURE thread_lock_check_trigger();

/*****************/

CREATE FUNCTION threads_prune_board_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	BEGIN
		DELETE FROM threads t
		USING last_reply_time l
		WHERE t.id = l.thread_id
		AND l.board_id = NEW.board_id
		AND l.sort_number > l.max_num_threads;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER threads_prune_board
	AFTER INSERT ON threads
	FOR EACH ROW
	EXECUTE PROCEDURE threads_prune_board_trigger();

/*****************/

CREATE FUNCTION comments_prune_board_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	BEGIN
		DELETE FROM threads t
		USING last_reply_time l
		WHERE t.id = l.thread_id
		AND l.board_id = (SELECT board_id FROM threads WHERE threads.id = NEW.reply_to)
		AND l.sort_number > l.max_num_threads;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER comments_prune_board
	AFTER INSERT ON comments
	FOR EACH ROW
	EXECUTE PROCEDURE comments_prune_board_trigger();

/*****************/

CREATE FUNCTION threads_post_num_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	DECLARE
		v_postnum INTEGER;
	BEGIN
		UPDATE board_postnum SET postnum = postnum + 1 WHERE board_id = NEW.board_id returning postnum INTO STRICT v_postnum;
		NEW.post_num := v_postnum;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER threads_update_postnum
	BEFORE INSERT ON threads
	FOR EACH ROW
	EXECUTE PROCEDURE threads_post_num_trigger();

/*****************/

CREATE FUNCTION comments_post_num_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	DECLARE
		v_postnum INTEGER;
	BEGIN
		UPDATE board_postnum SET postnum = postnum + 1 WHERE board_id = (SELECT board_id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE threads.id = NEW.reply_to) returning postnum INTO STRICT v_postnum;
		NEW.post_num := v_postnum;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER comments_update_postnum
	BEFORE INSERT ON comments
	FOR EACH ROW
	EXECUTE PROCEDURE comments_post_num_trigger();

/*****************/

CREATE FUNCTION make_board_postnum() RETURNS trigger
	LANGUAGE plpgsql AS $$
	BEGIN
		INSERT INTO board_postnum (board_id) VALUES (NEW.id);
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER update_postnum_table
	AFTER INSERT ON boards
	FOR EACH ROW
	EXECUTE PROCEDURE make_board_postnum();

/* CREATE VIEW posts_view AS */
/* 	SELECT boards.name AS board_name, posts.post_num, posts.reply_to, posts.image, posts.time, posts.comment */
/* 	FROM posts */
/* 	INNER JOIN boards */
/* 	ON posts.board_id = boards.id; */

/* CREATE VIEW original_posts AS */
/* 	SELECT * FROM posts_view */
/* 	WHERE reply_to IS NULL; */

/* CREATE VIEW replies AS */
/* 	SELECT * FROM posts_view */
/* 	WHERE reply_to IS NOT NULL; */

COMMIT;
