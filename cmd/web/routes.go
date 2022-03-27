package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/swissarmybox/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	requireAuth := dynamicMiddleware.Append(app.requireAuthentication)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	mux.Get("/snippet/create", requireAuth.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", requireAuth.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", requireAuth.ThenFunc(app.logoutUser))
	mux.Get("/user/profile", requireAuth.ThenFunc(app.userProfile))
	mux.Get("/user/change-password", requireAuth.ThenFunc(app.changePasswordForm))
	mux.Post("/user/change-password", requireAuth.ThenFunc(app.changePassword))

	mux.Get("/ping", http.HandlerFunc(ping))

	// For embed
	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Get("/static/", fileServer)

	// For non embedded
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
