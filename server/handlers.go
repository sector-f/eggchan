package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/sector-f/eggchan"
)

func isPretty(r *http.Request) bool {
	prettyVals := r.URL.Query()["pretty"]
	for _, value := range prettyVals {
		if value == "true" || value == "yes" || value == "please" {
			return true
		}
	}

	return false
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotFound, isPretty(r), "Not found")
}

// GET /
func (s *HttpServer) index(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	tabWriter := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)

	for _, route := range s.routes {
		message := fmt.Sprintf("%v\t%v\t%v\n", route.Method, route.Pattern, route.Description)
		tabWriter.Write([]byte(message))
	}
	tabWriter.Flush()

	w.Write(buf.Bytes())
}

// GET /categories
func (e *HttpServer) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := e.BoardService.ListCategories()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, isPretty(r), err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, isPretty(r), categories)
}

// GET /categories/{category}
func (e *HttpServer) showCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["category"]

	boards, err := e.BoardService.ShowCategory(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Invalid category")
		return
	}

	respondWithJSON(w, http.StatusOK, isPretty(r), boards)
}

// GET /boards/{board}
func (e *HttpServer) showBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["board"]

	boardReply, err := e.BoardService.ShowBoardReply(name)
	switch err.(type) {
	case nil:
		respondWithJSON(w, http.StatusOK, isPretty(r), boardReply)
	case eggchan.BoardNotFoundError:
		respondWithError(w, http.StatusNotFound, isPretty(r), err.Error())
	case eggchan.DatabaseError:
		respondWithError(w, http.StatusInternalServerError, isPretty(r), err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, isPretty(r), "Unknown error")
	}
}

// GET /boards/{board}/{thread}
func (e *HttpServer) showThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardName := vars["board"]
	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Invalid thread ID")
		return
	}

	threadReply, err := e.BoardService.ShowThreadReply(boardName, thread)
	switch err.(type) {
	case nil:
		respondWithJSON(w, http.StatusOK, isPretty(r), threadReply)
	case eggchan.BoardNotFoundError:
		respondWithError(w, http.StatusNotFound, isPretty(r), err.Error())
	case eggchan.ThreadNotFoundError:
		respondWithError(w, http.StatusNotFound, isPretty(r), err.Error())
	case eggchan.DatabaseError:
		respondWithError(w, http.StatusInternalServerError, isPretty(r), err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, isPretty(r), "Unknown error")
	}
}

// GET /boards
func (e *HttpServer) getBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := e.BoardService.ListBoards()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, isPretty(r), err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, isPretty(r), boards)
}

// POST /boards/{board}
func (e *HttpServer) postThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, isPretty(r), "Comment exceeds length limit")
		return
	}

	r.ParseMultipartForm(32 << 20)
	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Comment cannot be empty")
		return
	}

	author := r.FormValue("author")
	if author == "" {
		author = "Anonymous"
	}

	subject := r.FormValue("subject")

	post_num, err := e.BoardService.MakeThread(board, comment, author, subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Error creating thread")
		return
	}

	respondWithJSON(w, http.StatusCreated, isPretty(r), eggchan.PostThreadResponse{post_num})
}

// POST /boards/{board}/{thread}
func (e *HttpServer) postReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, isPretty(r), "Invalid thread ID")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, isPretty(r), "Comment exceeds length limit")
		return
	}

	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Comment cannot be empty")
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	if author == "" {
		author = "Anonymous"
	}

	post_num, err := e.BoardService.MakeComment(board, thread, comment, author)
	if err != nil {
		if err.(*pq.Error).Message == "Thread has reached post limit" {
			respondWithError(w, http.StatusForbidden, isPretty(r), "Thread has reached post limit")
		} else {
			respondWithError(w, http.StatusInternalServerError, isPretty(r), "Error creating post")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, isPretty(r), eggchan.PostCommentResponse{thread, post_num})
}

// DELETE /boards/{board}/threads/{thread}
func (e *HttpServer) deleteThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Invalid thread ID")
		return
	}

	deleted_count, err := e.AdminService.DeleteThread(board, thread)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, isPretty(r), "Could not delete thread")
		return
	}

	switch deleted_count {
	case 0:
		respondWithJSON(w, http.StatusNotFound, isPretty(r), SuccessMessage{"Thread not found"})
	case 1:
		respondWithJSON(w, http.StatusOK, isPretty(r), SuccessMessage{"Thread deleted"})
	default:
		respondWithJSON(w, http.StatusInternalServerError, isPretty(r), SuccessMessage{"Multiple threads were deleted--this is probably an error"})
	}
}

// DELETE /boards/{board}/comments/{comment}
func (e *HttpServer) deleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["comment"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Invalid comment ID")
		return
	}

	deleted_count, err := e.AdminService.DeleteComment(board, thread)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, isPretty(r), "Could not delete comment")
		return
	}

	switch deleted_count {
	case 0:
		respondWithJSON(w, http.StatusNotFound, isPretty(r), SuccessMessage{"Comment not found"})
	case 1:
		respondWithJSON(w, http.StatusOK, isPretty(r), SuccessMessage{"Comment deleted"})
	default:
		respondWithJSON(w, http.StatusInternalServerError, isPretty(r), SuccessMessage{"Multiple comments were deleted--this is probably an error"})
	}
}

// POST /new/boards/{board}
func (s *HttpServer) createBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, isPretty(r), "Length limit exceeded")
		return
	}

	description := strings.TrimSpace(r.FormValue("description"))
	if description == "" {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Description cannot be empty")
		return
	}

	category := strings.TrimSpace(r.FormValue("category"))
	if category == "" {
		respondWithError(w, http.StatusBadRequest, isPretty(r), "Category cannot be empty")
		return
	}

	err = s.AdminService.AddBoard(board, description, category)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, isPretty(r), err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, isPretty(r), SuccessMessage{"Board created"})
}
