package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
	"github.com/shtayeb/rssfeed/http/session"
	"github.com/shtayeb/rssfeed/internal"
	"github.com/shtayeb/rssfeed/internal/database"
	"github.com/shtayeb/rssfeed/internal/models"
	"github.com/shtayeb/rssfeed/views"
	"github.com/shtayeb/rssfeed/views/emails"
)

type UserRegisterParams struct {
	Name                 string
	Username             string
	Email                string
	Password             string
	PasswordConfirmation string
	Errors               map[string]string
}

func HandlerChangePassword(w http.ResponseWriter, r *http.Request) {
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
	// log.Printf("current Password:%v ", currentPassword)
	// log.Printf("hashedCurrentPassword: %v ", hashedCurrentPassword)
	// log.Printf("Auth User password: %v ", user.Password)

	if checkPasswordHash(currentPassword, user.Password) {
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
	err = internal.DB.ChangeUserPassword(
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
	// return something for htmx
}

func ResetPasswordView(w http.ResponseWriter, r *http.Request) {
	// REST: /reset-password/{token}
	token := chi.URLParam(r, "token")
	// verify the token
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err := verifyToken(token, *internal.Config)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		log.Println("Failed to parse the token", err)
		return
	}

	views.ResetPassword(token).Render(r.Context(), w)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// REST: /reset-password
	ctx := r.Context()
	r.ParseForm()
	token := r.PostFormValue("token")
	password := r.PostFormValue("new_password")
	passwordConfirmation := r.PostFormValue("new_password_confirmation")

	if token == "" || password == "" || passwordConfirmation == "" {
		msgs := []map[string]string{
			{
				"msg_type": "error",
				"msg":      "All fields are necessary!",
			},
		}
		RenderWithMsg(views.ResetPassword(token), w, ctx, msgs)
		return
	}

	if password != passwordConfirmation {
		RenderWithMsg(views.ResetPassword(token), w, ctx, []map[string]string{
			{
				"msg_type": "error",
				"msg":      "Password and Password Confirmation does not match!",
			},
		})
		return
	}

	jwtUserInfo, err := verifyToken(token, *internal.Config)
	if err != nil {
		RenderWithMsg(views.ResetPassword(token), w, ctx, []map[string]string{
			{
				"msg_type": "error",
				"msg":      "Something is wrong with your token",
			},
		})
		return
	}

	// Hash the new password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("failed to hash passwrod: %v", err)
		RenderWithMsg(views.ResetPassword(token), w, ctx, []map[string]string{
			{
				"msg_type": "error",
				"msg":      "Something went wrong, Please try again",
			},
		})
		return
	}

	// Update the user in db
	err = internal.DB.ChangeUserPassword(
		ctx,
		database.ChangeUserPasswordParams{Password: hashedPassword, ID: jwtUserInfo.user_id},
	)

	if err != nil {
		log.Println("Something went wrong,Please try again", err)
		return
	}

	// send the user to the login page
	ctx = context.WithValue(r.Context(), "msgs", []map[string]string{
		{
			"msg_type": "success",
			"msg":      "Your password has been reset, you can login with your new password",
		},
	})

	http.Redirect(w, r.WithContext(ctx), "/login", http.StatusSeeOther)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
	}

	email := r.PostFormValue("email")
	if email == "" {
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "Invalid Email"},
		}
		ctx := context.WithValue(r.Context(), "msgs", msgs)
		views.ForgotPassword().Render(ctx, w)
		return
	}

	log.Printf("email from the form: %v", email)

	// Get the user with the provided email
	user, err := internal.DB.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		msgs := []map[string]string{
			{"msg_type": "error", "msg": "No user with this email found!"},
		}
		ctx := context.WithValue(r.Context(), "msgs", msgs)
		views.ForgotPassword().Render(ctx, w)
		return
	}

	tokenString, err := createToken(user, *internal.Config)
	if err != nil {
		log.Println("Failed to create token:", err)
	}
	log.Printf("==== ==== token: %v", tokenString)

	// TODO: Email sending logic and configuration here
	mr := MailRequest{
		to:      []string{user.Email},
		subject: "RSS Feed: Reset Password Link",
	}

	// send the rest link to the user email address
	ok, err := SendEmail(
		ctx,
		emails.ResetPasswordRequestMail(
			user.Name,
			fmt.Sprintf("%v/reset-password/%v", internal.Config.APP_URL, tokenString),
		),
		*internal.Config,
		mr,
	)
	if !ok && err != nil {
		log.Printf("failed to send email: %v", err)
		msgs := []map[string]string{
			{
				"msg_type": "error",
				"msg":      "Could not send email, Please try again!",
			},
		}

		RenderWithMsg(views.ForgotPassword(), w, ctx, msgs)
		return
	}

	msgs := []map[string]string{
		{
			"msg_type": "success",
			"msg":      "Password reset link has been sent to your email address.",
		},
	}

	RenderWithMsg(views.ForgotPassword(), w, ctx, msgs)
}

func ForgotPasswordView(w http.ResponseWriter, r *http.Request) {
	views.ForgotPassword().Render(r.Context(), w)
}

func HandlerRegisterView(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.Register(map[string]string{}).Render(r.Context(), w)
}

func HandlerLoginView(w http.ResponseWriter, r *http.Request) {
	contextUser := r.Context().Value("user")
	if contextUser != nil {
		// Getout you are already loggedin
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	views.Login().Render(r.Context(), w)
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
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
	user, err := internal.DB.GetUserByEmailOrUsername(r.Context(), emailOrUsername)
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

func HandlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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
	_, err = internal.DB.CreateUser(r.Context(), database.CreateUserParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  params.Username,
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Println("Here test", err)
		// RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
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

func HandlerUsersGet(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)
	RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func HandlerLogout(w http.ResponseWriter, r *http.Request) {
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
