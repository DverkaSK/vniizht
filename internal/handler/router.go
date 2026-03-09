package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"vniizht/internal/middleware"
	"vniizht/internal/model"
	"vniizht/internal/service"
)

func NewRouter(h *Handlers, auth *service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Authenticate(auth))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Auth.Login)
		r.With(middleware.RequireAuth).Post("/logout", h.Auth.Logout)
	})

	r.Get("/search", h.Search.Search)
	r.Get("/categories", h.Admin.ListCategories)
	r.Get("/tags", h.Admin.ListTags)

	r.Route("/questions", func(r chi.Router) {
		r.Get("/", h.Questions.List)
		r.With(middleware.RequireAuth).Post("/", h.Questions.Create)
		r.Route("/{questionID}", func(r chi.Router) {
			r.Get("/", h.Questions.Get)
			r.With(middleware.RequireAuth).Patch("/", h.Questions.Update)
			r.With(middleware.RequireAuth).Patch("/close", h.Questions.Close)
			r.Get("/answers", h.Answers.List)
			r.With(middleware.RequireAuth).Post("/answers", h.Answers.Create)
		})
	})

	r.Route("/answers/{answerID}", func(r chi.Router) {
		r.With(middleware.RequireAuth).Patch("/", h.Answers.Update)
		r.With(middleware.RequireAuth).Delete("/", h.Answers.Delete)
		r.With(middleware.RequireRole(model.RoleSpecialist, model.RoleAdmin)).Patch("/verify", h.Answers.Verify)
		r.With(middleware.RequireRole(model.RoleSpecialist, model.RoleAdmin)).Delete("/verify", h.Answers.Unverify)
		r.With(middleware.RequireAuth).Post("/vote", h.Answers.CreateVote)
		r.With(middleware.RequireAuth).Patch("/vote", h.Answers.UpdateVote)
		r.With(middleware.RequireAuth).Delete("/vote", h.Answers.DeleteVote)
		r.Get("/comments", h.Comments.List)
		r.With(middleware.RequireAuth).Post("/comments", h.Comments.Create)
	})

	r.Route("/comments/{commentID}", func(r chi.Router) {
		r.With(middleware.RequireAuth).Patch("/", h.Comments.Update)
		r.With(middleware.RequireAuth).Delete("/", h.Comments.Delete)
	})

	r.Route("/attachments", func(r chi.Router) {
		r.Get("/", h.Attachments.List)
		r.With(middleware.RequireAuth).Post("/", h.Attachments.Upload)
		r.Get("/{attachmentID}", h.Attachments.Get)
		r.With(middleware.RequireAuth).Delete("/{attachmentID}", h.Attachments.Delete)
	})

	r.With(middleware.RequireAuth).Get("/users/me", h.Users.Me)
	r.Get("/users/{userID}", h.Users.GetPublic)

	adminOnly := middleware.RequireRole(model.RoleAdmin)
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminOnly)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.Admin.ListUsers)
			r.Post("/", h.Admin.CreateUser)
			r.Patch("/{userID}/role", h.Admin.ChangeRole)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", h.Admin.ListCategories)
			r.Post("/", h.Admin.CreateCategory)
			r.Patch("/{categoryID}", h.Admin.UpdateCategory)
			r.Delete("/{categoryID}", h.Admin.DeleteCategory)
		})
		r.Route("/tags", func(r chi.Router) {
			r.Get("/", h.Admin.ListTags)
			r.Post("/", h.Admin.CreateTag)
			r.Patch("/{tagID}", h.Admin.UpdateTag)
			r.Delete("/{tagID}", h.Admin.DeleteTag)
		})
	})

	return r
}
