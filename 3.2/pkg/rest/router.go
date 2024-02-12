package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter(auth *Auth, cntrl *UserController) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", cntrl.Login).Methods("POST")
	r.Handle("/user/{uuid}", auth.RegisteredUserJwtAuth(http.HandlerFunc(cntrl.View))).Methods("GET")
	r.Handle("/user/{uuid}", auth.AdminJwtAuth(http.HandlerFunc(cntrl.Edit))).Methods("PUT")
	r.Handle("/user/{uuid}", auth.AdminJwtAuth(http.HandlerFunc(cntrl.Delete))).Methods("DELETE")
	r.Handle("/user/{uuid}/vote", auth.RegisteredUserJwtAuth(http.HandlerFunc(cntrl.VoteProfile))).Methods("POST")

	return r
}
