package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hisamcode/acis/internal/data"
	"github.com/hisamcode/acis/internal/validator"
)

type projectForm struct {
	Name                string `form:"name"`
	Detail              string `form:"detail"`
	validator.Validator `form:"-"`
}

func (app *application) projects(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = projectForm{}
	app.render(w, http.StatusOK, LayoutStandard, "projects.html", data)
}
func (app *application) latestProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := app.DB.Project.LatestByUserID(app.userID)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Projects = projects
	app.render(w, http.StatusOK, LayoutPartials, "latest-projects.html", data)
}

type transactionForm struct {
	Nominal             float64 `form:"nominal"`
	Detail              string  `form:"detail"`
	WTSID               int8    `form:"wts_id"`
	EmojiID             string  `form:"emoji_id"`
	EmojiName           string  `form:"emoji_name"`
	Emoji               string  `form:"emoji"`
	validator.Validator `form:"-"`
}

func (app *application) project(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(r.PathValue("projectID"), 10, 64)
	if err != nil {
		app.renderServerError(w, err)
		return
	}
	project, err := app.DB.Project.Get(projectID)
	if err != nil {
		app.renderServerError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Project = *project
	data.Form = transactionForm{}
	app.render(w, http.StatusOK, LayoutProject, "home.html", data)
}

func (app *application) projectTransactionPost(w http.ResponseWriter, r *http.Request) {

	// validation
	form := transactionForm{}
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	transaction := data.Transaction{
		Nominal: form.Nominal,
		Detail:  form.Detail,
		WTSID:   1,
	}

	if len(form.EmojiID) > 0 {
		emoji := data.Emoji{}
		err = emoji.Decode(form.EmojiID)
		if err != nil {
			app.renderServerError(w, err)
			return
		}
		transaction.EmojiID = emoji
	} else {
		transaction.EmojiID = data.Emoji{
			ID:    "empty",
			Name:  "empty",
			Emoji: "empty",
		}
	}

	form.Validator = *validator.New()
	if data.ValidateTransaction(&form.Validator, &transaction); !form.Validator.Valid() {
		project, err := app.getProject(r)
		if err != nil {
			app.renderServerError(w, err)
			return
		}
		data := app.newTemplateData(r)
		data.Form = form
		data.Project = *project
		app.addHXTriggerAfterSettle(w, "validationCreateTransaction")
		app.render(w, http.StatusOK, LayoutProject, "home.html", data)
		return
	}

	project, err := app.getProject(r)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	transaction.ProjectID = project.ID
	transaction.CreatedBy = app.userID

	err = app.DB.Transaction.Insert(&transaction)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%d/home", project.ID), http.StatusSeeOther)
}

type emojiForm struct {
	ID    string `form:"emoji_id"`
	Name  string `form:"emoji_name"`
	Emoji string `form:"emoji"`
	validator.Validator
}

func (app *application) projectEmojiPut(w http.ResponseWriter, r *http.Request) {
	project, err := app.getProject(r)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	form := emojiForm{}

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	emoji := data.Emoji{
		ID:    form.ID,
		Name:  form.Name,
		Emoji: form.Emoji,
	}

	form.Validator = *validator.New()
	if data.ValidateEmoji(&form.Validator, &emoji); !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.addHXReswap(w, HXSWAP_INNER)
		app.render(w, http.StatusUnprocessableEntity, LayoutPartials, "validation-error.html", data)
		return
	}

	err = project.UpdateEmoji(emoji)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	err = app.DB.Project.Update(project)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Project = *project
	app.render(w, http.StatusOK, LayoutPartials, "list-emojis.html", data)
}

func (app *application) projectEmojiPost(w http.ResponseWriter, r *http.Request) {
	project, err := app.getProject(r)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	form := emojiForm{}

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	emoji := data.CreateEmoji(form.Name, form.Emoji)

	form.Validator = *validator.New()
	if data.ValidateEmoji(&form.Validator, &emoji); !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, LayoutPartials, "validation-error.html", data)
		return
	}

	project.Emojis = append(project.Emojis, emoji)

	err = app.DB.Project.Update(project)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Project = *project
	app.render(w, http.StatusOK, LayoutPartials, "list-emojis.html", data)
}

func (app *application) projectEmojiDelete(w http.ResponseWriter, r *http.Request) {
	project, err := app.getProject(r)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	form := emojiForm{}
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	emoji := data.Emoji{
		ID: form.ID,
	}

	err = project.DeleteEmoji(emoji)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	err = app.DB.Project.Update(project)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Project = *project
	app.render(w, http.StatusOK, LayoutPartials, "list-emojis.html", data)
}

type projectTransactionForm struct {
	transactionForm
	// projectForm for update on page setting
	Project projectForm
	Emoji   emojiForm
}

func (app *application) projectSetting(w http.ResponseWriter, r *http.Request) {
	project, err := app.getProject(r)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	pf := projectTransactionForm{
		Project: projectForm{
			Name:   project.Name,
			Detail: project.Detail,
		},
		Emoji: emojiForm{},
	}

	d := app.newTemplateData(r)
	d.Project = *project
	d.Form = pf
	app.render(w, http.StatusOK, LayoutProject, "setting.html", d)
}

func (app *application) projectSettingPut(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(r.PathValue("projectID"), 10, 64)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	var form projectForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	project, err := app.DB.Project.Get(projectID)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	if len(form.Name) > 0 {
		project.Name = form.Name
	}

	if len(form.Detail) > 0 {
		project.Detail = form.Detail
	}

	form.Validator = *validator.New()
	if data.ValidateProject(&form.Validator, project); !form.Valid() {
		data := app.newTemplateData(r)
		data.Project = *project
		data.Form = form
		app.render(w, http.StatusOK, LayoutProject, "setting.html", data)
		return
	}

	err = app.DB.Project.Update(project)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.renderEditConflict(w, err)
		default:
			app.renderServerError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%d/settings", projectID), http.StatusSeeOther)
}

func (app *application) projectPost(w http.ResponseWriter, r *http.Request) {
	var form projectForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.renderServerError(w, err)
		return
	}

	project := data.Project{
		Name:   form.Name,
		Detail: form.Detail,
		UserID: app.userID,
	}

	form.Validator = *validator.New()
	if data.ValidateProject(&form.Validator, &project); !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.addHXTriggerAfterSettle(w, "validationCreateProject")
		app.render(w, http.StatusOK, LayoutStandard, "projects.html", data)
		return
	}

	id, err := app.DB.Project.Insert(&project)
	if err != nil {
		app.renderServerError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%d/home", id), http.StatusSeeOther)
}
