package handlers

import (
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"vk_db_project/app/models"
	"vk_db_project/app/uscases"
)

type UserHandler struct {
	userLogic uscases.IUserUsecase
}

func NewUserHandler(uLogic uscases.UserUsecaseImpl) UserHandler {
	return UserHandler{userLogic: uLogic}
}

func (User UserHandler) GetStatus(rwContext echo.Context) error {
	return rwContext.JSON(http.StatusOK, User.userLogic.GetServerStatus())
}

func (User UserHandler) Clear(rwContext echo.Context) error {
	User.userLogic.Clear()
	return rwContext.NoContent(http.StatusOK)
}

func (User UserHandler) CreateUser(rwContext echo.Context) error {
	nickname := rwContext.Param("nickname")
	newUserData := new(models.UserModel)
	rwContext.Bind(newUserData)
	newUserData.Nickname = nickname
	answer, err := User.userLogic.CreateUser(*newUserData)
	if err != nil {
		return rwContext.JSON(http.StatusConflict, answer)
	}

	return rwContext.JSON(http.StatusCreated, answer)

}

func (User UserHandler) GetUser(rwContext echo.Context) error {
	nickname := rwContext.Param("nickname")
	userData, err := User.userLogic.GetUser(nickname)
	if err != nil {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't find user by nickname: " + nickname})
	}

	return rwContext.JSON(http.StatusOK, userData)
}

func (User UserHandler) UpdateUser(rwContext echo.Context) error {
	nickname := rwContext.Param("nickname")
	newUserData := new(models.UserModel)
	rwContext.Bind(newUserData)
	newUserData.Nickname = nickname

	answer, err := User.userLogic.UpdateUserData(*newUserData)
	if err == pgx.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message: "Can't find user by nickname: " + nickname})
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, &models.Error{Message: "This email is already registered by user: " + nickname})
	}

	return rwContext.JSON(http.StatusOK, answer)
}

func (User UserHandler) SetupHandlers(server *echo.Echo) {
	server.POST("/api/user/:nickname/create", User.CreateUser)
	server.GET("/api/user/:nickname/profile", User.GetUser)
	server.POST("/api/user/:nickname/profile", User.UpdateUser)
	server.GET("/api/service/status", User.GetStatus)
	server.POST("/api/service/clear", User.Clear)
}

