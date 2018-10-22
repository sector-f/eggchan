CREATE TABLE IF NOT EXISTS categories (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS boards (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	category INTEGER REFERENCES categories
);

/* Table board_postnum keeps track of the highest post number on each board */
CREATE TABLE IF NOT EXISTS board_postnum (
	board_id INTEGER PRIMARY KEY REFERENCES boards,
	postnum INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS images (
	id SERIAL PRIMARY KEY,
	filepath TEXT NOT NULL,
	thumbpath TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS threads (
	id SERIAL PRIMARY KEY,
	subject TEXT,
	author TEXT DEFAULT 'Anonymous',
	post_num INTEGER,
	board_id INTEGER REFERENCES boards,
	image INTEGER REFERENCES images,
	time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	comment TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
	id SERIAL PRIMARY KEY,
	author TEXT DEFAULT 'Anonymous',
	post_num INTEGER,
	reply_to INTEGER REFERENCES threads,
	image INTEGER REFERENCES images,
	time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	comment TEXT NOT NULL
);

CREATE VIEW comments_with_board_info AS
	SELECT b.id AS board_id, b.name AS board_name, c.id, c.post_num, c.reply_to, c.image, c.time, c.comment
	FROM comments c
	INNER JOIN boards b ON b.id = (SELECT board_id FROM threads INNER JOIN boards ON threads.board_id = boards.id);


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
