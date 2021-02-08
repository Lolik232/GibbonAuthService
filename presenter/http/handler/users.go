package handler

import (
	errors "auth-server/internal/app/errors/types"
	"auth-server/internal/app/model"
	"auth-server/internal/app/service"
	"auth-server/internal/app/store"
	"auth-server/internal/app/utils/emailsender"
	"auth-server/internal/app/utils/files"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

type UserHandler struct {
	Handler
	serviceManager *service.Manager
	emailSender    *emailsender.IEmailSender
}

func NewUserHandler(manager *service.Manager, sender *emailsender.IEmailSender) *UserHandler {
	return &UserHandler{
		Handler:        Handler{},
		serviceManager: manager,
		emailSender:    sender,
	}
}

func parseUserParams(params []string) *store.UserFields {
	fields := &store.UserFields{
		UserName:     false,
		Email:        false,
		CreatedAt:    false,
		UserInfo:     false,
		UserSessions: false,
		UserRoles:    false,
	}
	if params == nil {
		return fields
	}
	for _, param := range params {
		switch param {
		case store.ParamUserName:
			fields.UserName = true
		case store.ParamEmail:
			fields.Email = true
		case store.ParamCreatedAt:
			fields.CreatedAt = true
		case store.ParamUserInfo:
			fields.UserInfo = true
		case store.ParamUserRoles:
			fields.UserRoles = true
		case store.ParamUserSessions:
			fields.UserSessions = true
		}
	}
	return fields
}

func (u UserHandler) ConfigureRoutes(router *mux.Router) {

	users := router.PathPrefix("/users").Subrouter()
	//get user
	users.HandleFunc("/get/id/{id}", u.getUserByID()).Methods(http.MethodGet)
	users.HandleFunc("/get/username/{username}", u.getUserByName()).Methods(http.MethodGet)
	//register
	users.HandleFunc("/register", u.register()).Methods(http.MethodPost)
	users.HandleFunc("/email/confirm/{token}", u.confirmEmail()).Methods(http.MethodGet)
}

func (u UserHandler) getUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accepted client. Method: getUserByID, handler: user.")
		vars := mux.Vars(r)
		userId := vars["id"]
		params := parseUserParams(r.Form["params"])
		user, err := (*u.serviceManager.User).FindUserByID(r.Context(), userId, params)
		if err != nil {
			u.error(w, r, err)
		}
		u.respondJson(w, r, http.StatusOK, user)
	}
}
func (u UserHandler) getUserByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accepted client. Method: getUserByName, handler: user.")
		vars := mux.Vars(r)
		username := vars["username"]
		params := parseUserParams(r.Form["params"])
		user, err := (*u.serviceManager.User).FindUserByID(r.Context(), username, params)
		if err != nil {
			u.error(w, r, err)
		}
		u.respondJson(w, r, http.StatusOK, user)
	}
}

func (u UserHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accepted client. Method: register, handler: user.")
		user := &model.User{
			ID:             "",
			UserName:       "",
			Email:          "",
			EmailConfirmed: false,
			Password:       "",
			PasswordHash:   "",
			UserInfo:       nil,
			UserSessions:   nil,
			CreatedAt:      nil,
			Roles:          nil,
		}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			err = errors.ErrInvalidArgument.New("Invalid user data.")
			u.error(w, r, err)
		}
		token, err := (*u.serviceManager.User).Registration(r.Context(), user)
		if err != nil {
			u.error(w, r, err)
		}
		domain := viper.GetString("domain")
		tmpl, err := template.ParseFiles(filesPath.EmailConfTemplate)

		if err != nil {
			log.Println("Faili proebal, debil!")
			err = (*u.serviceManager.User).DeleteByName(r.Context(), user.UserName)
			if err != nil {
				log.Println("Polnii pizdec! Err in delete user.")
			}
			u.error(w, r, err)
		}

		data := struct {
			Email string
			Link  string
		}{
			Email: user.Email,
			Link:  fmt.Sprintf("%s/users/email/comfirm/%s", domain, token),
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, data)
		if err != nil {
			log.Println(err)
			u.error(w, r, err)
		}
		msg := buf.String()
		err = (*u.emailSender).Send(r.Context(), "Confirmation email", user.Email, "text/html", msg)
		if err != nil {
			log.Println(err)
			err = (*u.serviceManager.User).DeleteByName(r.Context(), user.UserName)
			if err != nil {
				log.Println("Polnii pizdec! Err in delete user.")
			}
			u.error(w, r, err)
		}
	}
}

func (u UserHandler) confirmEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		userID := r.FormValue("id")
		err := (*u.serviceManager.User).ConfirmEmail(r.Context(), userID, token)
		u.error(w, r, err)
	}
}
