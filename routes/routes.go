package routes

import (
	"net/http"
	"todo-auth/handler"
	"todo-auth/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func apiroutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.Caller)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", handler.Add)
		r.Get("/", handler.List)
		r.Put("/", handler.Update)
		r.Delete("/", handler.Delete)
	})
	return r

}
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/logout", handler.Logout)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method Not Found", http.StatusMethodNotAllowed)
		})
		r.Mount("/", apiroutes())
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}
