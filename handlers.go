package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotFound, "Not found")
}

func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategoriesFromDB(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (a *App) showCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["category"]

	boards, err := showCategoryFromDB(a.DB, name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}

func (a *App) showBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["board"]

	posts, err := showBoardFromDB(a.DB, name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (a *App) showThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid thread ID")
		return
	}

	posts, err := showThreadFromDB(a.DB, board, thread)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (a *App) getBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := getBoardsFromDB(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}
