package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/guregu/null.v3"
	"time"
)

type category struct {
	Name string `json:"name"`
}

func getCategoriesFromDB(db *sql.DB) ([]category, error) {
	rows, err := db.Query("SELECT name FROM categories")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []category{}
	for rows.Next() {
		var c category
		if err := rows.Scan(&c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

type board struct {
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Category    null.String `json:"category"`
}

func getBoardsFromDB(db *sql.DB) ([]board, error) {
	rows, err := db.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func showCategoryFromDB(db *sql.DB, name string) ([]board, error) {
	rows, err := db.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

type thread struct {
	PostNum         int       `json:"post_num"`
	Time            time.Time `json:"post_time"`
	NumReplies      int       `json:"num_replies"`
	LatestReply     null.Time `json:"latest_reply_time"`
	Comment         string    `json:"comment"`
	SortLatestReply time.Time `json:"-"`
}

type post struct {
	PostNum int       `json:"post_num"`
	Time    time.Time `json:"time"`
	Comment string    `json:"comment"`
}

func showBoardFromDB(db *sql.DB, name string) ([]thread, error) {
	rows, err := db.Query(
		`SELECT
			threads.post_num,
			threads.time,
			(SELECT COUNT(*) FROM comments WHERE comments.reply_to = threads.id) as num_replies,
			MAX(comments.time) AS latest_reply,
			threads.comment,
			CASE
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END AS sort_latest_reply
		FROM threads
		LEFT JOIN comments ON threads.id = comments.reply_to
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		GROUP BY threads.id
		ORDER BY sort_latest_reply DESC`,
		name,
	)
	// 	GROUP BY original_posts.board_name, original_posts.time, original_posts.post_num, original_posts.comment

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	threads := []thread{}
	for rows.Next() {
		var t thread
		if err := rows.Scan(&t.PostNum, &t.Time, &t.NumReplies, &t.LatestReply, &t.Comment, &t.SortLatestReply); err != nil {
			fmt.Println(err)
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func showThreadFromDB(db *sql.DB, board string, thread int) ([]post, error) {
	rows, err := db.Query(
		`SELECT threads.post_num, threads.time, threads.comment
		FROM threads
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE boards.name = $1
		AND threads.post_num = $2
		UNION
		SELECT comments.post_num, comments.time, comments.comment
		FROM comments
		INNER JOIN threads ON comments.reply_to = (SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2)
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		ORDER BY post_num ASC`,
		board,
		thread,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.PostNum, &p.Time, &p.Comment); err != nil {
			fmt.Println(err)
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func makeThreadInDB(db *sql.DB, board string, comment string) (int, error) {
	rows, err := db.Query(
		`INSERT INTO threads (board_id, comment)
		VALUES((SELECT id FROM boards WHERE name = $1), $2)
		RETURNING post_num`,
		board,
		comment,
	)

	if err != nil {
		fmt.Println(err)
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
		`INSERT INTO comments (reply_to, comment)
		VALUES(
			(SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2),
			$3
		)
		RETURNING post_num`,
		board,
		thread,
		comment,
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			fmt.Println(err)
			return 0, err
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func checkIsOp(db *sql.DB, board string, thread int) (bool, error) {
	rows, err := db.Query(
		`SELECT post_num
		FROM threads
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2`,
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
