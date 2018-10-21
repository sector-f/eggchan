package main

import (
	"database/sql"
	"gopkg.in/guregu/null.v3"
	"time"
)

type category struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

func (c *category) getCategoryFromDB(db *sql.DB) error {
	return db.QueryRow("SELECT id, name FROM categories WHERE id=$1", c.ID).Scan(&c.ID, &c.Name)
}

// func getCategories(db *sql.DB, start, count int) ([]category, error) {
func getCategoriesFromDB(db *sql.DB) ([]category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []category{}
	for rows.Next() {
		var c category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

type board struct {
	ID          int         `json:"-"`
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Category    null.String `json:"category"`
}

func getBoardsFromDB(db *sql.DB) ([]board, error) {
	rows, err := db.Query("SELECT boards.id, boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func showCategoryFromDB(db *sql.DB, name string) ([]board, error) {
	rows, err := db.Query("SELECT boards.id, boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

// NB: threads and posts are stored in the same table

type thread struct {
	PostNum         int       `json:"post_num"`
	Time            time.Time `json:"post_time"`
	LatestReply     null.Time `json:"latest_reply_time"`
	Comment         string    `json:"comment"`
	SortLatestReply time.Time `json:"-"`
}

type post struct {
	PostNum int       `json:"post_num"`
	ReplyTo null.Int  `json:"reply_to"`
	Time    time.Time `json:"time"`
	Comment string    `json:"comment"`
}

func showBoardFromDB(db *sql.DB, name string) ([]thread, error) {
	rows, err := db.Query(
		`SELECT original_posts.post_num, original_posts.time, MAX(replies.time) AS latest_reply, original_posts.comment,
			CASE
				WHEN MAX(replies.time) IS NOT NULL THEN MAX(replies.time)
				ELSE MAX(original_posts.time)
			END AS sort_latest_reply
		FROM original_posts
		LEFT JOIN replies ON original_posts.post_num = replies.reply_to
		WHERE original_posts.board_name = $1
		GROUP BY original_posts.board_name, original_posts.time, original_posts.post_num, original_posts.comment
		ORDER BY sort_latest_reply DESC;`,
		name,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []thread{}
	for rows.Next() {
		var t thread
		if err := rows.Scan(&t.PostNum, &t.Time, &t.LatestReply, &t.Comment, &t.SortLatestReply); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func showThreadFromDB(db *sql.DB, board string, thread int) ([]post, error) {
	rows, err := db.Query(
		`SELECT DISTINCT posts.post_num, posts.reply_to, posts.time, posts.comment
		FROM posts
		INNER JOIN boards ON (SELECT id FROM boards WHERE name = $1 LIMIT 1) = posts.board_id
		WHERE posts.post_num = $2 OR posts.reply_to = $2
		ORDER BY posts.post_num, posts.time ASC`,
		board,
		thread,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.PostNum, &p.ReplyTo, &p.Time, &p.Comment); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func makeThreadInDB(db *sql.DB, board string, comment string) (int, error) {
	rows, err := db.Query(
		`INSERT INTO posts (board_id, comment)
		VALUES((SELECT id FROM boards WHERE name = $1), $2)
		RETURNING post_num`,
		board,
		comment,
	)

	if err != nil {
		return 0, err
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, err
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func makePostInDB(db *sql.DB, board string, thread int, comment string) (int, error) {
	rows, err := db.Query(
		`INSERT INTO posts (board_id, reply_to, comment)
		VALUES((SELECT id FROM boards WHERE name = $1), $2, $3)
		RETURNING post_num`,
		board,
		thread,
		comment,
	)

	if err != nil {
		return 0, err
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, err
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func checkIsOp(db *sql.DB, board string, thread int) (bool, error) {
	rows, err := db.Query(
		`SELECT post_num
		FROM original_posts
		WHERE board_name = $1
		AND post_num = $2`,
		board,
		thread,
	)

	if err != nil {
		return false, err
	}

	posts := []post{}
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.PostNum); err != nil {
			return false, err
		}
		posts = append(posts, p)
	}

	if len(posts) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
