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
	board, err := r.Service.ShowBoardDesc(obj.name)
	if err != nil {
		return "", err
	}
	return board.Name, nil
}
func (r *boardResolver) Description(ctx context.Context, obj *Board) (*string, error) {
	board, err := r.Service.ShowBoardDesc(obj.name)
	if err != nil {
		return nil, err
	}

	if board.Description.Valid {
		return &board.Description.String, nil
	} else {
		return nil, nil
	}
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

type postResolver struct{ *Resolver }

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
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

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

type threadResolver struct{ *Resolver }

func (r *threadResolver) PostNum(ctx context.Context, obj *Thread) (int, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return 0, err
	}

	return thread.PostNum, nil
}
func (r *threadResolver) Subject(ctx context.Context, obj *Thread) (*string, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return nil, err
	}

	return &thread.Subject.String, nil
}
func (r *threadResolver) Author(ctx context.Context, obj *Thread) (string, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return "", err
	}

	return thread.Author, nil
}
func (r *threadResolver) Time(ctx context.Context, obj *Thread) (time.Time, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return time.Now(), err
	}

	return thread.Time, nil
}
func (r *threadResolver) NumReplies(ctx context.Context, obj *Thread) (int, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return 0, err
	}

	return thread.NumReplies, nil
}
func (r *threadResolver) LatestReplyTime(ctx context.Context, obj *Thread) (time.Time, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return time.Now(), err
	}

	return thread.SortLatestReply, nil
}
func (r *threadResolver) Comment(ctx context.Context, obj *Thread) (string, error) {
	thread, err := r.Service.ShowThreadOP(obj.board, obj.postNum)
	if err != nil {
		return "", err
	}

	return thread.Comment, nil
}
func (r *threadResolver) Posts(ctx context.Context, obj *Thread) ([]*Post, error) {
	panic("not implemented")
}
