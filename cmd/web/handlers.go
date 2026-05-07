package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/thedevscott/trug/internal/models"
	"github.com/thedevscott/trug/internal/validator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type accountPasswordUpdateForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

// Transaction represents a single financial movement.
type transactionCreateForm struct {
	Title               string `form:"title"`
	IsIncome            bool   `form:"is_income"`
	Amount              string `form:"amount"`
	Category            string `form:"category"`
	Description         string `form:"description"`
	Date                string `form:"date" time_format:"2006-01-02"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// transactionCreate A handler for creating a snippet
func (app *application) transactionCreate(w http.ResponseWriter, r *http.Request) {
	userID := app.getCurrentUsersID(w, r)
	transactions, err := app.transactions.Latest(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// 30 days ago
	start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// present day
	end := time.Now().Format("2006-01-02")

	stats, err := app.transactions.GetUserStats(userID, start, end)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Transactions = transactions
	data.Stats = stats

	data.Form = transactionCreateForm{
		Date: time.Now().Format("2006-01-02"),
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) transactionCreatePost(w http.ResponseWriter, r *http.Request) {
	var form transactionCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Category), "category", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	transactionDate, err := time.Parse("2006-01-02", form.Date)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	amount, err := strconv.ParseFloat(form.Amount, 64)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	amountInCents := amount * 100

	isChecked := r.FormValue("is_income") != ""
	form.IsIncome = isChecked

	userID := app.getCurrentUsersID(w, r)
	_, err = app.transactions.Insert(userID, form.Title, form.IsIncome, int64(amountInCents), form.Category, form.Description, transactionDate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	caser := cases.Title(language.English)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("%s Transaction recorded!", caser.String(form.Title)))

	http.Redirect(w, r, "/transaction/create", http.StatusSeeOther)
}

func (app *application) transactionView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	userID := app.getCurrentUsersID(w, r)
	transaction, err := app.transactions.Get(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Form = transactionCreateForm{
		Title:       transaction.Title,
		IsIncome:    transaction.IsIncome,
		Amount:      fmt.Sprintf("%.2f", float64(transaction.Amount/100)),
		Category:    transaction.Category,
		Description: transaction.Description,
		Date:        transaction.TransactionDate.Format("2006-01-02"),
	}

	data.Transaction = transaction

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) transactionUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var form transactionCreateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Category), "category", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, fmt.Sprintf("/transaction/view/%d", id), data)
		return
	}

	transactionDate, err := time.Parse("2006-01-02", form.Date)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	amount, err := strconv.ParseFloat(form.Amount, 64)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	amountInCents := amount * 100

	isChecked := r.FormValue("is_income") != ""
	form.IsIncome = isChecked

	userID := app.getCurrentUsersID(w, r)
	_, err = app.transactions.Update(userID, id, form.Title, form.IsIncome, int64(amountInCents), form.Category, form.Description, transactionDate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	caser := cases.Title(language.English)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("%s Transaction updated!", caser.String(form.Title)))

	http.Redirect(w, r, fmt.Sprintf("/transaction/view/%d", id), http.StatusSeeOther)
}

func (app *application) transactionDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	userID := app.getCurrentUsersID(w, r)
	num, err := app.transactions.Delete(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Deleted %d entry", num))

	http.Redirect(w, r, "/transaction/create", http.StatusSeeOther)
}

// about the about page "/about" for the app
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about.tmpl.html", data)
}

// accountView route for the logged in users profile
func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account.tmpl.html", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	path := app.sessionManager.PopString(r.Context(), "afterLoginRedirect")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/transaction/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// de-authorize and 'logout' user
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}

	app.render(w, r, http.StatusOK, "password.tmpl.html", data)
}

func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl.html", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.PasswordUpdate(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated!")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) getCurrentUsersID(w http.ResponseWriter, r *http.Request) int {
	currentUsersID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if currentUsersID == 0 {
		app.serverError(w, r, errors.New("failded to verify user"))
	}
	return currentUsersID
}
