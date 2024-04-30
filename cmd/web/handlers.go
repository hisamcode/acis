package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hisamcode/acis/internal/data"
	"github.com/hisamcode/acis/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, LayoutStandard, "home.html", templateData{})
}
func (app *application) transactionCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, LayoutClean, "transaction-create.html", templateData{})
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

type SignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	RepeatPassword      string `form:"repeat_password"`
	validator.Validator `form:"-"`
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = SignupForm{}
	app.render(w, LayoutClean, "signup.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {

	var form SignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user := data.User{
		Name:      form.Name,
		Email:     form.Email,
		Activated: false,
	}

	err = user.Password.Set(form.Password)
	if err != nil {
		// TODO: error server
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Validator = *validator.New()
	form.Validator.Check(validator.Equal(form.Password, form.RepeatPassword), "repeatPassword", "Repeat password not same")
	if data.ValidateUser(&form.Validator, &user); !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, LayoutClean, "signup.html", data)
		return
	}

	err = app.DB.User.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			// TODO: duplicate error
			form.AddFieldError("email", "a user with this email adress already exists")
			data := app.newTemplateData()
			data.Form = form
			app.render(w, LayoutClean, "signup.html", data)
		default:
			// TODO server errro
			app.logger.Error(err.Error())
		}
		return
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.html", user)
		if err != nil {
			app.logger.Error(err.Error())
			fmt.Fprintf(w, "sending email failed: %v", err)
			return
		}
	})

	fmt.Fprintf(w, "create a new user...")

}
