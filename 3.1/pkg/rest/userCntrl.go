package rest

import (
	"net/http"
	"strconv"
	"usr_mngmnt/pkg/usecase"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// Rest controllers.
type userCntrl struct {
	ucase usecase.UserUsecaseInterface
}

func NewUserCntrl(ucase usecase.UserUsecaseInterface) *userCntrl {
	return &userCntrl{ucase: ucase}
}

// Create user profile.
func (cntrl *userCntrl) ViewAll(c echo.Context) error {
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed offset")
	}

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed limit")
	}

	users := cntrl.ucase.ViewUsers(offset, limit)

	return c.JSON(http.StatusOK, users)
}

// Create user profile.
func (cntrl *userCntrl) View(c echo.Context) error {
	uuid := c.Param("uuid")
	user := cntrl.ucase.View(uuid)
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found")
	}

	return c.JSON(http.StatusCreated, user)
}

// Create user profile.
func (cntrl *userCntrl) Create(c echo.Context) error {
	// Bind request data to struct.
	req := new(usecase.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	// Validate request.
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = cntrl.ucase.Create(*req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, req)
}

// Update user profile.
func (cntrl *userCntrl) Update(c echo.Context) error {
	// Bind request data to struct.
	req := new(usecase.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	// Validate request.
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = cntrl.ucase.Update(*req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, req)
}

func (cntrl *userCntrl) BasicAuthMiddleware(email, pass string, _ echo.Context) (bool, error) {
	if cntrl.ucase.IsAuthenticated(email, pass) {
		return true, nil
	}

	return false, nil
}
