//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"
	time "time"

	"github.com/sector-f/eggchan"
)

type Resolver struct {
	Service eggchan.BoardService
}

func (r *Resolver) Board() BoardResolver {
	return &boardResolver{r}
}
func (r *Resolver) Category() CategoryResolver {
	return &categoryResolver{r}
}
func (r *Resolver) Post() PostResolver {
	return &postResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Thread() ThreadResolver {
	return &threadResolver{r}
}

type boardResolver struct{ *Resolver }

func (r *boardResolver) Name(ctx context.Context, obj *Board) (string, error) {
	return obj.name, nil
}
func (r *boardResolver) Description(ctx context.Context, obj *Board) (*string, error) {
	return &obj.description, nil
}
func (r *boardResolver) Category(ctx context.Context, obj *Board) (*string, error) {
	board, err := r.Service.ShowBoardDesc(obj.name)
	if err != nil {
		return nil, err
	}

	if board.Category.Valid {
		return &board.Category.String, nil
	} else {
		return nil, nil
	}
}
func (r *boardResolver) Threads(ctx context.Context, obj *Board) ([]*Thread, error) {
	threadReply := []*Thread{}

	threads, err := r.Service.ShowBoard(obj.name)
	if err != nil {
		return nil, err
	}

	for _, thread := range threads {
		posts := []int{}

		postsInThread, err := r.Service.ShowThread(thread.Board, thread.PostNum)
		if err != nil {
			return nil, err
		}

		for _, post := range postsInThread {
			posts = append(posts, post.PostNum)
		}

		newThread := Thread{
			board:           thread.Board,
			postNum:         thread.PostNum,
			subject:         thread.Subject.String,
			author:          thread.Author,
			numReplies:      thread.NumReplies,
			latestReplyTime: thread.SortLatestReply,
			comment:         thread.Comment,
			posts:           posts,
		}

		threadReply = append(threadReply, &newThread)
	}

	return threadReply, nil
}

type categoryResolver struct{ *Resolver }

func (r *categoryResolver) Name(ctx context.Context, obj *Category) (string, error) {
	return obj.name, nil
}
func (r *categoryResolver) Boards(ctx context.Context, obj *Category) ([]*Board, error) {
	boards, err := r.Service.ListBoards()
	if err != nil {
		return nil, err
	}

	var boardReply []*Board
	for _, board := range boards {
		threads := []int{}

		threadsInBoard, err := r.Service.ShowBoard(board.Name)
		if err != nil {
			return nil, err
		}

		for _, thread := range threadsInBoard {
			threads = append(threads, thread.PostNum)
		}

		newBoard := Board{
			name:        board.Name,
			description: board.Description.String,
			category:    board.Category.String,
			threads:     threads,
		}

		if board.Category.String == obj.name {
			boardReply = append(boardReply, &newBoard)
		}
	}

	return boardReply, nil
}

type postResolver struct{ *Resolver }

func (r *postResolver) Board(ctx context.Context, obj *Post) (string, error) {
	panic("not implemented")
}
func (r *postResolver) PostNum(ctx context.Context, obj *Post) (int, error) {
	panic("not implemented")
}
func (r *postResolver) Author(ctx context.Context, obj *Post) (string, error) {
	panic("not implemented")
}
func (r *postResolver) Time(ctx context.Context, obj *Post) (time.Time, error) {
	panic("not implemented")
}
func (r *postResolver) Comment(ctx context.Context, obj *Post) (string, error) {
	return obj.comment, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Categories(ctx context.Context, name *string) ([]*Category, error) {
	categoryReply := []*Category{}
	if name == nil {
		categories, err := r.Service.ListCategories()
		if err != nil {
			return nil, err
		}

		for _, category := range categories {
			boards, err := r.Service.ShowCategory(category.Name)
			if err != nil {
				return nil, err
			}

			boardList := []string{}
			for _, board := range boards {
				boardList = append(boardList, board.Name)
			}

			newCategory := Category{
				name:   category.Name,
				boards: boardList,
			}

			categoryReply = append(categoryReply, &newCategory)
		}
	} else {
		boards, err := r.Service.ShowCategory(*name)
		if err != nil {
			return nil, err
		}

		boardList := []string{}
		for _, board := range boards {
			boardList = append(boardList, board.Name)
		}

		newCategory := Category{
			name:   *name,
			boards: boardList,
		}

		categoryReply = append(categoryReply, &newCategory)
	}

	return categoryReply, nil
}
func (r *queryResolver) Boards(ctx context.Context, name *string) ([]*Board, error) {
	boardReply := []*Board{}

	boards, err := r.Service.ListBoards()
	if err != nil {
		return nil, err
	}

	for _, board := range boards {
		if name == nil || board.Name == *name {
			threads := []int{}

			threadsInBoard, err := r.Service.ShowBoard(board.Name)
			if err != nil {
				return nil, err
			}

			for _, thread := range threadsInBoard {
				threads = append(threads, thread.PostNum)
			}

			newBoard := Board{
				name:        board.Name,
				description: board.Description.String,
				category:    board.Category.String,
				threads:     threads,
			}

			boardReply = append(boardReply, &newBoard)
		}
	}

	return boardReply, nil
}
func (r *queryResolver) Thread(ctx context.Context, board string, id int) (Thread, error) {
	thread, err := r.Service.ShowThreadOP(board, id)
	if err != nil {
		return Thread{}, err
	}

	posts, err := r.Service.ShowThread(board, id)
	if err != nil {
		return Thread{}, err
	}

	postNums := []int{}
	for _, post := range posts {
		postNums = append(postNums, post.PostNum)
	}

	threadReply := Thread{
		board: thread.Board,
		postNum: thread.PostNum,
		subject: thread.Subject.String,
		author: thread.Author,
		time: thread.Time,
		numReplies: thread.NumReplies,
		latestReplyTime: thread.SortLatestReply,
		comment: thread.Comment,
		posts: postNums,
	}

	return threadReply, nil
}

type threadResolver struct{ *Resolver }

func (r *threadResolver) PostNum(ctx context.Context, obj *Thread) (int, error) {
	return obj.postNum, nil
}
func (r *threadResolver) Subject(ctx context.Context, obj *Thread) (*string, error) {
	return &obj.subject, nil
}
func (r *threadResolver) Author(ctx context.Context, obj *Thread) (string, error) {
	return obj.author, nil
}
func (r *threadResolver) Time(ctx context.Context, obj *Thread) (time.Time, error) {
	return obj.time, nil
}
func (r *threadResolver) NumReplies(ctx context.Context, obj *Thread) (int, error) {
	return obj.numReplies, nil
}
func (r *threadResolver) LatestReplyTime(ctx context.Context, obj *Thread) (time.Time, error) {
	return obj.latestReplyTime, nil
}
func (r *threadResolver) Comment(ctx context.Context, obj *Thread) (string, error) {
	return obj.comment, nil
}
func (r *threadResolver) Posts(ctx context.Context, obj *Thread) ([]*Post, error) {
	posts, err := r.Service.ShowThread(obj.board, obj.postNum)
	if err != nil {
		return nil, err
	}

	postReply := []*Post{}
	for _, post := range posts {
		newPost := Post{
			board: post.Board,
			postNum: post.PostNum,
			author: post.Author,
			time: post.Time,
			comment: post.Comment,
		}

		postReply = append(postReply, &newPost)
	}

	return postReply, nil
}
