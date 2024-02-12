package rest

import (
	"encoding/json"
	"fmt"
	"jwt/pkg/cache"
	"jwt/pkg/infra"
	"jwt/pkg/model"
	"jwt/pkg/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type UserController struct {
	ucase           usecase.UserUsecaseInterface
	votecase        *usecase.Vote
	validator       *validator.Validate
	rspndr          ResponderInterface
	log             infra.LoggerInterface
	jwtService      JwtServiceInterface
	cache           cache.CacheInterface
	cacheTimeoutSec int
}

func NewUserController(ucase *usecase.User, votecase *usecase.Vote, validator *validator.Validate, rspndr *responder, logger infra.LoggerInterface, jwtService JwtServiceInterface, cache cache.CacheInterface, cacheTimeoutSec int) *UserController {
	return &UserController{
		ucase:           ucase,
		votecase:        votecase,
		validator:       validator,
		rspndr:          rspndr,
		log:             logger,
		jwtService:      jwtService,
		cache:           cache,
		cacheTimeoutSec: cacheTimeoutSec,
	}
}

func (cntrl *UserController) View(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	var user model.UserVotes
	err := cntrl.cache.Get(uuid, user)
	// Cached value does not exist, create it.
	if err != nil {
		user, err := cntrl.ucase.View(uuid)
		if err != nil {
			cntrl.log.Warn("Usecase:" + err.Error())
			cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = cntrl.cache.Set(uuid, user, time.Second*time.Duration(cntrl.cacheTimeoutSec))
		if err != nil {
			cntrl.log.Warn("Cache:" + err.Error())
			cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	cntrl.rspndr.Success(w, user)
}

func (cntrl *UserController) Edit(w http.ResponseWriter, r *http.Request) {
	var userRequest usecase.UserRequest

	// Bind request to user struct.
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		cntrl.log.Warn("jsonDecorder: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request.
	err = cntrl.validator.Struct(userRequest)
	if err != nil {
		cntrl.log.Warn("Validator: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Init usecase.
	err = cntrl.ucase.Update(&userRequest)
	if err != nil {
		cntrl.log.Warn("Usecase:" + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cntrl.rspndr.Success(w, userRequest)
}

func (cntrl *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err := cntrl.ucase.Delete(uuid)
	if err != nil {
		cntrl.log.Warn("Usecase: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cntrl.rspndr.Success(w, "User deleted.")
}

type VoteRequest struct {
	FromUser string `validate:"required" json:"from_user"`
	ToUser   string `validate:"required" json:"to_user"`
	Vote     string `validate:"required" json:"vote"`
}

func (cntrl *UserController) VoteProfile(w http.ResponseWriter, r *http.Request) {

	// Bind request to rqst struct.
	rqst := VoteRequest{}
	err := json.NewDecoder(r.Body).Decode(&rqst)
	if err != nil {
		cntrl.log.Warn("jsonDecorder: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	rqst.ToUser = vars["uuid"]

	// Validate request.
	err = cntrl.validator.Struct(rqst)
	if err != nil {
		cntrl.log.Warn("Validator: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rqst.Vote != strconv.Itoa(usecase.UPVOTE) && rqst.Vote != strconv.Itoa(usecase.DOWNVOTE) {
		errMsg := fmt.Sprintf("vote payload must be either %d or %d", usecase.UPVOTE, usecase.DOWNVOTE)
		cntrl.log.Warn(errMsg)
		cntrl.rspndr.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Call usecase.
	vote, err := strconv.Atoi(rqst.Vote)
	if err != nil {
		cntrl.log.Warn("Vote: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch vote {
	case usecase.UPVOTE:
		err = cntrl.votecase.Upvote(rqst.FromUser, rqst.ToUser)
	case usecase.DOWNVOTE:
		err = cntrl.votecase.Downvote(rqst.FromUser, rqst.ToUser)
	}

	if err != nil {
		cntrl.log.Warn("Usecase: " + err.Error())
		cntrl.rspndr.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cntrl.rspndr.Success(w, "Vote has been cast.")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func (cntrl *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		cntrl.log.Warn("Error decoding credentials")
		cntrl.rspndr.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := cntrl.ucase.Auth(loginRequest.Email, loginRequest.Password)
	if err != nil {
		cntrl.log.Warn(err)
		cntrl.rspndr.Error(w, "Login credentials incorrect. Cannot log you in.", http.StatusUnauthorized)
		return
	}

	token, err := cntrl.jwtService.GenerateToken(*user.Email, *user.Role)
	if err != nil {
		cntrl.log.Warn(err)
		cntrl.rspndr.Error(w, "Something went wrong on our side, sorry!", http.StatusInternalServerError)
		return
	}

	cntrl.rspndr.Success(w, LoginResponse{Token: token})

	cntrl.log.Info("User logged in and received a token")
}
