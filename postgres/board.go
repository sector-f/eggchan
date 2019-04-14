package postgres

import (
	"database/sql"

	"github.com/sector-f/eggchan"
)

func (s *EggchanService) ShowBoardReply(board string) (eggchan.BoardReply, error) {
	return eggchan.BoardReply{}, eggchan.UnimplementedError{}
}

func (s *EggchanService) ShowThreadReply(board string, id int) (eggchan.ThreadReply, error) {
	return eggchan.ThreadReply{}, eggchan.UnimplementedError{}
}

func (s *EggchanService) ListCategories() ([]eggchan.Category, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}
	defer tx.Commit()

	catRows, err := tx.Query("SELECT name FROM categories ORDER BY name ASC")

	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	categories := []eggchan.Category{}
	for catRows.Next() {
		var c eggchan.Category
		if err := catRows.Scan(&c.Name); err != nil {
			return nil, eggchan.DatabaseError{}
		}

		categories = append(categories, c)
	}
	catRows.Close()

	for i, category := range categories {
		boardRows, err := tx.Query("SELECT b.name, b.description, $1::text FROM boards b INNER JOIN categories c ON c.id = b.category WHERE c.name = $1::text", category.Name)
		if err != nil {
			return nil, eggchan.DatabaseError{}
		}

		boards := []eggchan.Board{}
		for boardRows.Next() {
			var b eggchan.Board
			if err := boardRows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
				return nil, eggchan.DatabaseError{}
			}
			boards = append(boards, b)
		}
		boardRows.Close()

		categories[i].Boards = boards
	}

	return categories, nil
}

func (s *EggchanService) ListBoards() ([]eggchan.Board, error) {
	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id ORDER BY boards.name ASC")

	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, eggchan.DatabaseError{}
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowCategory(name string) ([]eggchan.Board, error) {
	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1 ORDER BY boards.name ASC", name)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, eggchan.DatabaseError{}
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowBoard(name string) ([]eggchan.Thread, error) {
	rows, err := s.DB.Query(
		`SELECT
			$1::text,
			threads.post_num,
			threads.subject,
			threads.author,
			threads.time,
			(SELECT COUNT(*) FROM comments WHERE comments.reply_to = threads.id) AS num_replies,
			CASE
				WHEN MAX(comments.time) IS NOT NULL AND COUNT(*) >= (SELECT bump_limit FROM boards WHERE name = $1)  THEN (SELECT comments.time FROM comments OFFSET (SELECT bump_limit FROM boards WHERE name = $1) LIMIT 1)
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END AS sort_latest_reply,
			threads.comment
		FROM threads
		LEFT JOIN comments ON threads.id = comments.reply_to
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		GROUP BY threads.id
		ORDER BY sort_latest_reply DESC`,
		name,
	)

	threads := []eggchan.Thread{}

	if err != nil {
		return threads, eggchan.DatabaseError{}
	}
	defer rows.Close()

	for rows.Next() {
		var t eggchan.Thread
		if err := rows.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment); err != nil {
			return threads, eggchan.DatabaseError{}
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func (s *EggchanService) ShowThread(board string, thread_num int) ([]eggchan.Post, error) {
	c_rows, err := s.DB.Query(
		`SELECT comments.reply_to, comments.post_num, comments.author, comments.time, comments.comment
		FROM comments
		INNER JOIN threads ON comments.reply_to = threads.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND comments.reply_to = (SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2)
		ORDER BY post_num ASC`,
		board,
		thread_num,
	)

	posts := []eggchan.Post{}

	if err != nil {
		return posts, eggchan.DatabaseError{}
	}
	defer c_rows.Close()

	for c_rows.Next() {
		var p eggchan.Post
		if err := c_rows.Scan(&p.ReplyTo, &p.PostNum, &p.Author, &p.Time, &p.Comment); err != nil {
			return posts, eggchan.DatabaseError{}
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *EggchanService) MakeThread(board string, comment string, author string, subject string) (int, error) {
	rows, err := s.DB.Query(
		`INSERT INTO threads (board_id, comment, author, subject)
		VALUES(
			(SELECT id FROM boards WHERE name = $1),
			$2,
			$3,
			CASE WHEN $4 = '' THEN NULL ELSE $4 END
		)
		RETURNING post_num`,
		board,
		comment,
		author,
		subject,
	)

	if err != nil {
		return 0, eggchan.DatabaseError{}
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, eggchan.DatabaseError{}
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func (s *EggchanService) MakeComment(board string, thread int, comment string, author string) (int, error) {
	row := s.DB.QueryRow(
		`INSERT INTO comments (reply_to, comment, author)
		VALUES(
			(SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2),
			$3,
			$4
		)
		RETURNING post_num`,
		board,
		thread,
		comment,
		author,
	)

	var post_num int

	err := row.Scan(&post_num)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return 0, eggchan.NotFoundError{}
	default:
		return 0, eggchan.DatabaseError{}
	}

	return post_num, nil
}

func (s *EggchanService) checkIsOp(board string, thread int) (bool, error) {
	rows, err := s.DB.Query(
		`SELECT post_num
		FROM threads
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2`,
		board,
		thread,
	)

	if err != nil {
		return false, eggchan.DatabaseError{}
	}

	posts := []eggchan.Post{}
	for rows.Next() {
		var p eggchan.Post
		if err := rows.Scan(&p.PostNum); err != nil {
			return false, eggchan.DatabaseError{}
		}
		posts = append(posts, p)
	}

	if len(posts) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *EggchanService) ShowBoardDesc(board string) (eggchan.Board, error) {
	b_row := s.DB.QueryRow(
		`SELECT boards.name, boards.description, boards.category
		FROM boards
		WHERE boards.name = $1`,
		board,
	)

	var b eggchan.Board

	err := b_row.Scan(&b.Name, &b.Description, &b.Category)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return b, eggchan.NotFoundError{}
	default:
		return b, eggchan.DatabaseError{}
	}

	return b, nil
}

func (s *EggchanService) ShowThreadOP(board string, id int) (eggchan.Thread, error) {
	t_row := s.DB.QueryRow(
		`SELECT
			$1::text,
			threads.post_num,
			threads.subject,
			threads.author,
			threads.time,
			(SELECT COUNT(*) FROM comments WHERE comments.reply_to = threads.id) AS num_replies,
			CASE
				WHEN MAX(comments.time) IS NOT NULL AND COUNT(*) >= (SELECT bump_limit FROM boards WHERE name = $1)  THEN (SELECT comments.time FROM comments OFFSET (SELECT bump_limit FROM boards WHERE name = $1) LIMIT 1)
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END AS sort_latest_reply,
			threads.comment
		FROM threads
		LEFT JOIN comments ON threads.id = comments.reply_to
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2
		GROUP BY threads.id
		ORDER BY sort_latest_reply DESC`,
		board,
		id,
	)

	var t eggchan.Thread

	err := t_row.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return t, eggchan.NotFoundError{}
	default:
		return t, eggchan.DatabaseError{}
	}

	return t, nil
}
