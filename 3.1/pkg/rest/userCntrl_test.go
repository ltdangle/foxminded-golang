package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"usr_mngmnt/pkg/model"
	"usr_mngmnt/pkg/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRequestFail(t *testing.T) {
	// Init httptest.
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	// Init echo webserver.
	e := echo.New()
	c := e.NewContext(req, rec)

	cntrl := NewUserCntrl(usecase.NewUserUsecaseInterfaceStub())
	err := cntrl.Create(c)
	assert.Error(t, err)
}

func TestViewUserProfile(t *testing.T) {
	ucase := usecase.NewUserUsecaseInterfaceStub()
	ucase.SetView(func() *model.User {
		return model.NewUser()
	})

	cntrl := NewUserCntrl(ucase)

	// Init httptest.
	req := httptest.NewRequest(http.MethodGet, "/user/48493", nil)
	rec := httptest.NewRecorder()

	// Init echo webserver.
	e := echo.New()
	c := e.NewContext(req, rec)

	if assert.NoError(t, cntrl.View(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateUserProfile(t *testing.T) {
	cntrl := NewUserCntrl(usecase.NewUserUsecaseInterfaceStub())

	// Init httptest.
	requestJSON := `{"email":"jon@labstack.com", "password":"myveryveryverylongpassword", "first_name":"Jon", "last_name":"Snow", "nickname":"jsnow"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Init echo webserver.
	e := echo.New()
	c := e.NewContext(req, rec)

	if assert.NoError(t, cntrl.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestViewUsers(t *testing.T) {
	cntrl := NewUserCntrl(usecase.NewUserUsecaseInterfaceStub())

	// Init httptest.
	req := httptest.NewRequest(http.MethodGet, "/users/0/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Init echo webserver.
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:offset/:limit")
	c.SetParamNames("offset", "limit")
	c.SetParamValues("0", "10")

	if assert.NoError(t, cntrl.ViewAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUpdateProfile(t *testing.T) {
	cntrl := NewUserCntrl(usecase.NewUserUsecaseInterfaceStub())

	// Init httptest.
	requestJSON := `{"email":"jon@labstack.com", "password":"myveryveryverylongpassword", "first_name":"Jon", "last_name":"Snow", "nickname":"jsnow"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	req.SetBasicAuth("user", "pass")

	// Init echo webserver.
	e := echo.New()
	c := e.NewContext(req, rec)

	if assert.NoError(t, cntrl.Update(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
