package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shtayeb/rssfeed/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name                 string
		Username             string
		Email                string
		Password             string
		PasswordConfirmation string
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	// name := r.PostForm.Get("name")
	params := parameters{
		Name:     r.PostForm.Get("name"),
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("Email"),
		Password: r.PostForm.Get("password"),
	}

	_, err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  params.Username,
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
	})

	if err != nil {
		log.Println("Here test", err)
		// respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't create user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Added Successfully!"))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
