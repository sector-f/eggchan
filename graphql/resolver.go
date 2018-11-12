//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"
	time "time"
)

type Resolver struct{}

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
	panic("not implemented")
}
func (r *boardResolver) Description(ctx context.Context, obj *Board) (*string, error) {
	panic("not implemented")
}
func (r *boardResolver) Category(ctx context.Context, obj *Board) (*string, error) {
	panic("not implemented")
}
func (r *boardResolver) Threads(ctx context.Context, obj *Board) ([]*Thread, error) {
	panic("not implemented")
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

func (r *queryResolver) Boards(ctx context.Context) ([]*Board, error) {
	panic("not implemented")
}

type threadResolver struct{ *Resolver }

func (r *threadResolver) PostNum(ctx context.Context, obj *Thread) (int, error) {
	panic("not implemented")
}
func (r *threadResolver) Subject(ctx context.Context, obj *Thread) (*string, error) {
	panic("not implemented")
}
func (r *threadResolver) Author(ctx context.Context, obj *Thread) (string, error) {
	panic("not implemented")
}
func (r *threadResolver) Time(ctx context.Context, obj *Thread) (time.Time, error) {
	panic("not implemented")
}
func (r *threadResolver) NumReplies(ctx context.Context, obj *Thread) (int, error) {
	panic("not implemented")
}
func (r *threadResolver) LatestReplyTime(ctx context.Context, obj *Thread) (time.Time, error) {
	panic("not implemented")
}
func (r *threadResolver) Comment(ctx context.Context, obj *Thread) (string, error) {
	panic("not implemented")
}
func (r *threadResolver) Posts(ctx context.Context, obj *Thread) ([]*Post, error) {
	panic("not implemented")
}
