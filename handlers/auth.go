package handlers

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/angelofallars/htmx-go"
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

func (cfg *ApiConfig) HandlerChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get the user from session
	// Get compare the hashed passwords
	r.ParseForm()
	currentPassword := r.PostFormValue("current_password")
	newPassword := r.PostFormValue("new_password")
	newPasswordConfirmation := r.PostFormValue("new_password_confirmation")

	if currentPassword == "" || newPassword == "" || newPasswordConfirmation == "" {
		// Invalid data
		// return validation error
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Please fill all of the form !"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}
	// compare the password and password confirmation
	if newPassword != newPasswordConfirmation {
		// Password does not match
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Passwords does not match"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}

	user := ctx.Value("user").(database.User)

	hashedCurrentPassword, _ := hashPassword(currentPassword)
	log.Printf("current Password:%v ", currentPassword)
	log.Printf("hashedCurrentPassword: %v ", hashedCurrentPassword)
	log.Printf("Auth User password: %v ", user.Password)
	if user.Password != hashedCurrentPassword {
		// current password is wrong
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Your current password is wrong !"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}

	// Hash the nw password
	hashedNewPassword, err := hashPassword(newPassword)
	if err != nil {
		// Could not hash the new password
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Invalid new password. Please try something new"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}

	// Update the user in db
	err = cfg.DB.ChangeUserPassword(
		ctx,
		database.ChangeUserPasswordParams{Password: hashedNewPassword, ID: user.ID},
	)
	if err != nil {
		// Could not update the DB, try again
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Something went wrong please try again"},
		}
		htmx.NewResponse().
			RenderTempl(r.Context(), w, views.RenderMessages(msgs))
		return
	}
	// Invalidate the session - logout user
	err = session.SessionManager.Destroy(r.Context())
	if err != nil {
		log.Println("Failed to Destroy the session")
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
	// return something for htmx
}

func (cfg *ApiConfig) HandlerRegisterView(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.Register(map[string]string{}).Render(r.Context(), w)
}

func (cfg *ApiConfig) HandlerLoginView(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.Login().Render(r.Context(), w)
}

func hashPassword(password string) (string, error) {
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
	contextUser := r.Context().Value("user")
	if contextUser != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// If user is loggedin redirect them back with a mesasge
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	emailOrUsername := r.PostFormValue("email_or_username")
	password := r.PostFormValue("password")

	if emailOrUsername == "" || password == "" {
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Email and Password should not be empty !"},
		}
		ctx := context.WithValue(r.Context(), "msgs", msgs)
		views.Login().Render(ctx, w)
		return
	}

	// Ge the user by email or username
	user, err := cfg.DB.GetUserByEmailOrUsername(r.Context(), emailOrUsername)
	if err != nil {
		log.Println("Failed to get user from the Database", err)

		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Invalid Email or Password, Try Again! 1"},
		}
		ctx := context.WithValue(r.Context(), "msgs", msgs)

		views.Login().Render(ctx, w)
		return
	}

	// Get the password against hashed pass
	if !checkPasswordHash(password, user.Password) {
		// Email Invalid
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Invalid Email or Password, Try Again! 2"},
		}
		ctx := context.WithValue(r.Context(), "msgs", msgs)
		views.Login().Render(ctx, w)
		return
	}

	// Set the seesion ID in the cookie
	log.Printf("Login: loggedin user: %v ", user.ID)
	session.SessionManager.Put(r.Context(), "user_id", user.ID)

	user_id := session.SessionManager.GetInt(r.Context(), "user_id")

	log.Printf("Login: UserID in login: %v", user_id)
	//
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   (session.ID).String(),
	// 	Expires: time.Now().Add(120 * time.Second),
	// })
	// msg := []map[string]string{{"msg_type": "error", "msg": "Invalid Email or Password, Try Again!"}}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (cfg *ApiConfig) HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

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
		// msg := []map[string]string{}
		views.Register(params.Errors).Render(r.Context(), w)
		return
	}

	hashedPassword, err := hashPassword(params.Password)
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
		msgs := []map[string]string{{"msg_type": "error", "msg": "Could not create user"}}
		ctx := context.WithValue(r.Context(), "msgs", msgs)

		views.Register(params.Errors).Render(ctx, w)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	msgs := []map[string]string{{"msg_type": "success", "msg": "User created successfully !"}}
	ctx := context.WithValue(r.Context(), "msgs", msgs)
	views.Register(params.Errors).Render(ctx, w)
}

func (cfg *ApiConfig) HandlerUsersGet(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)
	respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func (cfg *ApiConfig) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	// Delete the session from the DB
	log.Println("LogoutHanlder")
	err := session.SessionManager.Destroy(r.Context())
	if err != nil {
		log.Println("Failed to Destroy the session")
		return
	}

	// context.WithValue(r.Context(), "user", nil)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
