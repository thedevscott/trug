package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/thedevscott/trug/ui"
)

// routes returns a ServeMux containing the applications routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// use embedded files
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// health check route
	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	// Register app routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /transaction/create", protected.ThenFunc(app.transactionCreate))
	mux.Handle("POST /transaction/create", protected.ThenFunc(app.transactionCreatePost))
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
