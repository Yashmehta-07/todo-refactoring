package routes

import (
	"todo/handler"
	middlewares "todo/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Route() chi.Router {

	r := chi.NewRouter()

	//middlewares
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(caller)
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/logout", handler.Logout)
		r.Route("/tasks", func(r chi.Router) {
			r.Use(middlewares.Caller)
			r.Get("/", handler.List)
			r.Post("/", handler.Add)
			r.Put("/", handler.Update)
			r.Delete("/", handler.Delete)
		})
	})
	// r.Get("/tasks", middlewares.Caller(handler.List))
	// r.Post("/tasks", middlewares.Caller(handler.Add))
	// r.Put("/tasks", middlewares.Caller(handler.Update))
	// r.Delete("/tasks", middlewares.Caller(handler.Delete))

	//user routes
	// r.Post("/login", handler.Login)
	// r.Post("/register", handler.Register)
	// r.Post("/logout", handler.Logout)

	return r
}
