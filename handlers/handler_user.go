package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
	"github.com/shtayeb/rssfeed/internal/session"
	"github.com/shtayeb/rssfeed/views"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterParams struct {
	Name                 string
	Username             string
	Email                string
	Password             string
	PasswordConfirmation string
	Errors               map[string]string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var rxEmail = regexp.MustCompile(".+@.+\\..+")

func (params *UserRegisterParams) Validate() bool {
	params.Errors = make(map[string]string)

	if !rxEmail.Match([]byte(params.Email)) {
		params.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(params.Username) == "" {
		params.Errors["Username"] = "Please enter a username"
	}
	if strings.TrimSpace(params.Name) == "" {
		params.Errors["Name"] = "Please enter a name"
	}

	// check for password and password confirmation
	if params.Password != params.PasswordConfirmation {
		params.Errors["Password"] = "Password does not match !"
	}

	return len(params.Errors) == 0
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// If user is loggedin redirect them back with a mesasge
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	emailOrUsername := r.PostFormValue("email_or_username")
	password := r.PostFormValue("password")

	if emailOrUsername == "" || password == "" {
		msg := []map[string]string{
			{"msg_type": "error", "msg": "Email and Password should not be empty !"},
		}
		views.Login(msg).Render(r.Context(), w)
		return
	}

	// Ge the user by email or username
	user, err := cfg.DB.GetUserByEmailOrUsername(r.Context(), emailOrUsername)
	if err != nil {
		log.Println("Failed to get user from the Database", err)

		msg := []map[string]string{
			{"msg_type": "error", "msg": "Invalid Email or Password, Try Again! 1"},
		}
		views.Login(msg).Render(r.Context(), w)
		return
	}

	// Get the password against hashed pass
	if !checkPasswordHash(password, user.Password) {
		// Email Invalid
		msg := []map[string]string{
			{"msg_type": "error", "msg": "Invalid Email or Password, Try Again! 2"},
		}
		views.Login(msg).Render(r.Context(), w)
		return
	}

	// ses, err := cfg.DB.CreateSession(r.Context(), database.CreateSessionParams{
	// 	ID:        uuid.New(),
	// 	UserID:    user.ID,
	// 	ExpiresAt: sql.NullTime{},
	// })
	// if err != nil {
	// 	log.Println("Failed to create seesion in the the Database", err)
	// 	msg := []map[string]string{
	// 		{"msg_type": "error", "msg": "Something went wrong. Please try again!"},
	// 	}
	// 	views.Login(msg).Render(r.Context(), w)
	//
	// 	return
	// }
	// Set the seesion ID in the cookie
	log.Printf("Login: loggedin user: %v ", user.ID)
	session.SessionManager.Put(r.Context(), "user_id", user.ID)

	// token, time, err := session.SessionManager.Commit(r.Context())
	// log.Printf("Login: Commit session: %v in the time: %v", token, time)

	user_id := session.SessionManager.GetInt(r.Context(), "user_id")

	log.Printf("Login: UserID in login: %v", user_id)
	//
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   (session.ID).String(),
	// 	Expires: time.Now().Add(120 * time.Second),
	// })
	// msg := []map[string]string{{"msg_type": "error", "msg": "Invalid Email or Password, Try Again!"}}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	// msg := []map[string]string{{"msg_type": "success", "msg": "Success !"}}
	// views.Login(msg).Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	params := UserRegisterParams{
		Name:                 r.PostFormValue("name"),
		Username:             r.PostFormValue("username"),
		Email:                r.PostFormValue("email"),
		Password:             r.PostFormValue("password"),
		PasswordConfirmation: r.PostFormValue("password_confirmation"),
	}

	// https://github.com/go-playground/validator
	if !params.Validate() {
		msg := []map[string]string{}
		views.Register(msg, params.Errors).Render(r.Context(), w)
		return
		// http.Redirect(w, r, "/register", http.StatusSeeOther)
		// https://blog.jetbrains.com/go/2022/11/08/build-a-blog-with-go-templates/#creating-the-routes
		// http.Error(w, fmt.Sprintf("Validation error"), http.StatusBadRequest)
		// return

	}

	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		log.Println("Failed to hash the password", err)
	}

	// Hash the password
	_, err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  params.Username,
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Println("Here test", err)
		// respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		// w.WriteHeader(http.StatusInternalServerError)
		msg := []map[string]string{{"msg_type": "error", "msg": "Could not create user"}}
		views.Register(msg, params.Errors).Render(r.Context(), w)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	msg := []map[string]string{{"msg_type": "success", "msg": "User created successfully !"}}
	views.Register(msg, params.Errors).Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func (cfg *ApiConfig) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	// Delete the session from the DB
	err := session.SessionManager.Destroy(r.Context())
	if err != nil {
		log.Println("Failed to Destroy the session")
		return
	}

	// context.WithValue(r.Context(), "user", nil)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
