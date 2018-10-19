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

CREATE TABLE IF NOT EXISTS board_postnum (
	board_id INTEGER PRIMARY KEY REFERENCES boards,
	postnum INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS images (
	id SERIAL PRIMARY KEY,
	filepath TEXT NOT NULL,
	thumbpath TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
	board_id INTEGER REFERENCES boards,
	post_num INTEGER NOT NULL,
	PRIMARY KEY(board_id, post_num),
	reply_to INTEGER,
	FOREIGN KEY (board_id, reply_to) REFERENCES posts (board_id, post_num),
	image INTEGER REFERENCES images,
	time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
	comment TEXT NOT NULL,
	CHECK (reply_to <> post_num)
);

CREATE FUNCTION post_num_trigger() RETURNS trigger
	LANGUAGE plpgsql AS $$
	DECLARE
		v_postnum INTEGER;
	BEGIN
		UPDATE board_postnum SET postnum = postnum + 1 WHERE board_id = NEW.board_id returning postnum INTO STRICT v_postnum;
		NEW.post_num := v_postnum;
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER update_postnum
	BEFORE INSERT ON posts
	FOR EACH ROW
	EXECUTE PROCEDURE post_num_trigger();

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

CREATE VIEW posts_view AS
	SELECT boards.name AS board_name, posts.post_num, posts.reply_to, posts.image, posts.time, posts.comment
	FROM posts
	INNER JOIN boards
	ON posts.board_id = boards.id;

CREATE VIEW original_posts AS
	SELECT * FROM posts_view
	WHERE reply_to IS NULL;

CREATE VIEW replies AS
	SELECT * FROM posts_view
	WHERE reply_to IS NOT NULL;

/* CREATE TABLE IF NOT EXISTS threads ( */
/* 	id INTEGER, */
/* 	board_id INTEGER REFERENCES boards(id), */
/* 	PRIMARY KEY(id, board_id) */
/* ); */

/* CREATE TABLE IF NOT EXISTS posts ( */
/* 	id INTEGER NOT NULL, */
/* 	thread_id INTEGER, */
/* 	board_id INTEGER, */
/* 	FOREIGN KEY (thread_id, board_id) REFERENCES threads (id, board_id), */
/* 	PRIMARY KEY(id, board_id), */
/* 	-- time TIMESTAMPTZ NOT NULL, */
/* 	comment TEXT NOT NULL, */
/* 	image INTEGER REFERENCES images(id) */
/* ); */
