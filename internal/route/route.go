package route

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/cors"
	"github.com/leminhson2398/todo-api/internal/db"
	"github.com/leminhson2398/todo-api/internal/frontend"
	"github.com/leminhson2398/todo-api/internal/graph"
)

type FrontendHandler struct {
	staticPath string
	indexPath  string
}

func (h FrontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f, err := frontend.Frontend.Open(path)
	if os.IsNotExist(err) || IsDir(f) {
		index, err := frontend.Frontend.Open("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, "index.html", time.Now(), index)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, path, time.Now(), f)
}

func IsDir(f http.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.IsDir()
}

type TodoHandler struct {
	repo   db.Repository
	jwtKey []byte
}

// NewRouter return main router for the whole api
func NewRouter(dbConnection *sqlx.DB, jwtKey []byte) (chi.Router, error) {

	cors_ := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	r.Use(cors_.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	repository := db.NewRepository(dbConnection)
	todoHandler := TodoHandler{*repository, jwtKey}

	var imgserver = http.FileServer(http.Dir("./uploads/"))
	r.Group(func(mux chi.Router) {
		mux.Mount("/auth", authResource{}.AuthGroup(todoHandler))
		mux.Handle("/__graphql", graph.NewPlaygroundHandler("/graphql"))
		mux.Mount("/uploads/", http.StripPrefix("/uploads/", imgserver))
	})

	auth := AuthenticationMiddleware{jwtKey}
	r.Group(func(mux chi.Router) {
		mux.Use(auth.Middleware)
		// mux.Post("/users/me/avatar", todoHandler.Insta)
		mux.Handle("/graphql", graph.NewHandler(*repository))
	})

	frontend := FrontendHandler{
		staticPath: "build",
		indexPath:  "index.html",
	}
	r.Handle("/*", frontend)

	return r, nil
}
