package route

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/leminhson2398/todo-api/internal/auth"
	"github.com/leminhson2398/todo-api/internal/db"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authResource struct{}

// LoginRequestData contains data for login request
type LoginRequestData struct {
	Email    string
	Password string
}

// LoginResponseData contains result for login request
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}

// LoginHandler login user
func (h *TodoHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestData LoginRequestData

	// decode login form data
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug("bad request body")
		return
	}

	// find user with user input email value
	user, err := h.repo.GetUserAccountByEmail(r.Context(), requestData.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"email": requestData.Email,
		}).Warn("user account not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate if user entered password match password in database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(requestData.Password))
	if err != nil {
		log.WithFields(log.Fields{
			"email": requestData.Email,
		}).Warn("password incorect for user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshCreatedAt := time.Now().UTC()
	refreshExpiresAt := refreshCreatedAt.AddDate(0, 0, 1)

	// create refresh token
	refreshTokenString, err := h.repo.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    user.UserID,
		CreatedAt: refreshCreatedAt,
		ExpiresAt: sql.NullTime{
			Time:  refreshExpiresAt,
			Valid: true,
		},
	})

	// create access token
	acessTokenString, err := auth.NewAccessToken(user.UserID.String(), auth.Unrestricted, user.RoleCode, h.jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// set refreshToken to cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshTokenString.ID.String(),
		Expires:  refreshExpiresAt,
		HttpOnly: true,
	})

	// response json data
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(LoginResponseData{
		AccessToken: acessTokenString,
	})
}

// LogoutResponseData contains logout result
type LogoutResponseData struct {
	Success string `json:"success"`
}

// LogoutHandler logout user
func (h *TodoHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshTokenID := uuid.MustParse(c.Value)
	err = h.repo.DeleteRefreshTokenByID(r.Context(), refreshTokenID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(LogoutResponseData{Success: "success"})
}

// func (h *TodoHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	_, err := r.repo
// }

// AuthGroup contains routes for authentication operations (login, logout, ...)
func (rs authResource) AuthGroup(todoHandler TodoHandler) chi.Router {
	router := chi.NewRouter()
	router.Post("/login", todoHandler.LoginHandler)
	router.Post("/logout", todoHandler.LogoutHandler)
	return router
}
