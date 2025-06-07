package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"go-jwt/token"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AppContext struct {
	DB *sql.DB
}

func RegisterRoutes(router *mux.Router, ac *AppContext) {
	log.Println("registering routes...")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", ac.LoginHandler).Methods("POST")
	apiRouter.HandleFunc("/signup", ac.SignUpHandler).Methods("POST")

	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(AuthMiddleware)

	protectedRouter.HandleFunc("/test", ac.ProtectedHandler).Methods("GET")
}

// TODO validate request bodies
// TODO set HTTP only cookie serverside
// TODO refresh token

func (ac *AppContext) SignUpHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var user userIn
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding message body: %s \n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ac.createUser(user)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				log.Printf("user already exists: %s \n", user.Email)
				writer.WriteHeader(http.StatusConflict)
				return
			}
		}
		log.Printf("failed to create user: %s \n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("user created")
	writer.WriteHeader(http.StatusCreated)
}

func (ac *AppContext) LoginHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var user userIn
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Printf("error decoding message body: %s \n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	usersDb, err := ac.getUserByEmail(user.Email)
	if err != nil {
		log.Printf("failed to retrieve user from db: %s \n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(usersDb) == 0 {
		log.Println("no users found with this email")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if usersDb[0].Email == user.Email && isItTheSamePw([]byte(usersDb[0].Password), []byte(user.Password)) {
		tokenString, err := token.CreateToken(user.Email)
		if err != nil {
			log.Printf("error creating token: %s \n", err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, tokenString)
		return
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(writer, "invalid credentials")
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		tokenString := req.Header.Get("Authorization")
		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		err := token.VerifyToken(tokenString)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, req)
	})
}

func (ac *AppContext) ProtectedHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	message := Message{
		Message: "Welcome to the protected area!",
	}
	messageJson, _ := json.Marshal(message)
	fmt.Fprint(writer, string(messageJson))
}
