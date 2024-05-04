package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hisamcode/acis/internal/data"
	"github.com/hisamcode/acis/internal/session"
	"github.com/hisamcode/acis/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, LayoutStandard, "home.html", data)
}
func (app *application) transactionCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, LayoutClean, "transaction-create.html", templateData{})
}
func (app *application) transactionPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) categoriesView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories page view")
}
func (app *application) categories(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories page list")
}
func (app *application) categoriesPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create categories")
}
func (app *application) categoriesEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories edit")
}
func (app *application) categoriesPut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories put")
}

type signupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	RepeatPassword      string `form:"repeat_password"`
	validator.Validator `form:"-"`
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, http.StatusOK, LayoutClean, "signup.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {

	var form signupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	user := data.User{
		Name:      form.Name,
		Email:     form.Email,
		Activated: false,
	}

	err = user.Password.Set(form.Password)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	form.Validator = *validator.New()
	form.Validator.Check(validator.Equal(form.Password, form.RepeatPassword), "repeatPassword", "Repeat password not same")
	if data.ValidateUser(&form.Validator, &user); !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, LayoutClean, "signup.html", data)
		return
	}

	err = app.DB.User.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			form.AddFieldError("email", "a user with this email adress already exists")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusOK, LayoutClean, "signup.html", data)
		default:
			app.renderServerError(w, err)
		}
		return
	}

	durationActivation, _ := time.ParseDuration(app.config.activationAccountDuration)
	token, err := app.DB.Token.New(user.ID, durationActivation, data.ScopeActivation)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
			"link":            app.config.host,
		}
		app.logger.Info("Send email activation", "email", user.Email)
		err = app.mailer.Send(user.Email, "user_welcome.html", data)
		if err != nil {
			app.logger.Error(err.Error())
			return
		}
	})

	http.Redirect(w, r, "/user/activated?is=check-email", http.StatusSeeOther)
}

func (app *application) activateAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: flash message after registered like thankyou for registered,
	// please check your email for activated your account

	p := r.URL.Query().Get("p")

	if p == "token" {
		token := r.URL.Query().Get("token")

		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			data := app.newTemplateData(r)
			data.Form = v
			data.TokenActivate = token
			app.render(w, http.StatusOK, LayoutClean, "activate-account.html", data)
			return
		}

		user, err := app.DB.User.GetForToken(data.ScopeActivation, token)
		if err != nil {
			app.logger.Error(err.Error())
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				v.AddFieldError("token", "token is invalid or expired activation token")
				data := app.newTemplateData(r)
				data.Form = v
				data.TokenActivate = token
				app.render(w, http.StatusOK, LayoutClean, "activate-account.html", data)
			default:
				app.renderServerError(w, err)
			}
			return
		}

		user.Activated = true

		err = app.DB.User.Update(user)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrEditConflict):
				app.renderEditConflict(w, err)
			default:
				app.renderServerError(w, err)
			}
			return
		}

		err = app.DB.Token.DeleteAllForUser(data.ScopeActivation, user.ID)
		if err != nil {
			app.renderServerError(w, err)
			return
		}

		// TODO: message activation success
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = struct {
		FieldErrors map[string]string
	}{
		FieldErrors: make(map[string]string),
	}
	app.render(w, http.StatusOK, LayoutClean, "activate-account.html", data)

}

type signinForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	RememberMe          bool   `form:"remember_me"`
	validator.Validator `form:"-"`
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	form := signinForm{}
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, http.StatusOK, LayoutClean, "signin.html", data)
}

func (app *application) signinPost(w http.ResponseWriter, r *http.Request) {

	form := signinForm{}
	form.Validator = *validator.New()

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
	}

	data.ValidateEmail(&form.Validator, form.Email)
	data.ValidatePasswordPlaintext(&form.Validator, form.Password)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, LayoutClean, "signin.html", data)
		return
	}

	user, err := app.DB.User.GetByEmail(form.Email)
	if err != nil {
		form.AddFieldError("email", "email not found")
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, LayoutClean, "signin.html", data)
		return
	}

	passwordOk, err := user.Password.Matches(form.Password)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	if !passwordOk {
		form.AddFieldError("password", "password not match")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, LayoutClean, "signin.html", data)
		return
	}

	if !user.Activated {
		form.AddFieldError("email", "User not activated")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, LayoutClean, "signin.html", data)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), session.SessionAuthenticatedUserID, user.ID)
	if form.RememberMe {
		app.sessionManager.RememberMe(r.Context(), true)
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), session.SessionAuthenticatedUserID)
	app.sessionManager.Put(r.Context(), session.SessionFlash, "You've been logged out succesfully")
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
