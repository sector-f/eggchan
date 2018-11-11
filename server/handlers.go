package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotFound, "Not found")
}

func (e *HttpServer) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := e.BoardService.ListCategories()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (e *HttpServer) showCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["category"]

	boards, err := e.BoardService.ShowCategory(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}

func (e *HttpServer) showBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["board"]

	posts, err := e.BoardService.ShowBoard(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (e *HttpServer) showThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid thread ID")
		return
	}

	posts, err := e.BoardService.ShowThread(board, thread)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid thread or board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (e *HttpServer) getBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := e.BoardService.ListBoards()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}

func (e *HttpServer) postThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, "Comment exceeds length limit")
		return
	}

	r.ParseMultipartForm(32 << 20)
	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	author := r.FormValue("author")
	if author == "" {
		author = "Anonymous"
	}

	subject := r.FormValue("subject")

	post_num, err := e.BoardService.MakeThread(board, comment, author, subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating thread")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"post_num": post_num})
}

func (e *HttpServer) postReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid thread ID")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, "Comment exceeds length limit")
		return
	}

	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	if author == "" {
		author = "Anonymous"
	}

	post_num, err := e.BoardService.MakeComment(board, thread, comment, author)
	if err != nil {
		if err.(*pq.Error).Message == "Thread has reached post limit" {
			respondWithError(w, http.StatusForbidden, "Thread has reached post limit")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error creating post")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"post_num": post_num})
}

func (e *HttpServer) deleteThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid thread ID")
		return
	}

	deleted_count, err := e.AdminService.DeleteThread(board, thread)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete thread")
		return
	}

	switch deleted_count {
	case 0:
		respondWithJSON(w, http.StatusNotFound, SuccessMessage{"Thread not found"})
	case 1:
		respondWithJSON(w, http.StatusOK, SuccessMessage{"Thread deleted"})
	default:
		respondWithJSON(w, http.StatusInternalServerError, SuccessMessage{"Multiple threads were deleted--this is probably an error"})
	}
}

func (e *HttpServer) deleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["comment"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	deleted_count, err := e.AdminService.DeleteComment(board, thread)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete comment")
		return
	}

	switch deleted_count {
	case 0:
		respondWithJSON(w, http.StatusNotFound, SuccessMessage{"Comment not found"})
	case 1:
		respondWithJSON(w, http.StatusOK, SuccessMessage{"Comment deleted"})
	default:
		respondWithJSON(w, http.StatusInternalServerError, SuccessMessage{"Multiple comments were deleted--this is probably an error"})
	}
}
