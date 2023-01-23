package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"vk_db_project/app/models"
	"vk_db_project/app/uscases"
)

type PostHandler struct {
	PostLogic uscases.IPostUsecase
}

func NewPostHandler(pLogic uscases.PostUsecaseImpl) PostHandler {
	return PostHandler{PostLogic: pLogic}
}

func (PostHandler PostHandler) GetPost(rwContext echo.Context) error {
	id, _ := strconv.Atoi(rwContext.Param("id"))
	related := rwContext.QueryParams()

	val := strings.Split(related.Get("related"), ",")

	allPostData, err := PostHandler.PostLogic.GetPostData(id, val)

	if err != nil {

		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find post by id: " + rwContext.Param("id")})
	}

	return rwContext.JSON(http.StatusOK, allPostData)
}

func (PostHandler PostHandler) UpdatePost(rwContext echo.Context) error {
	id, _ := strconv.ParseInt(rwContext.Param("id"), 10, 64)

	msg := new(models.Post)
	rwContext.Bind(msg)

	currentMsg, err := PostHandler.PostLogic.UpdatePost(id, msg.Message)
	if err != nil {

		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find post by id: " + rwContext.Param("id")})
	}

	return rwContext.JSON(http.StatusOK, currentMsg)
}

func (PostHandler PostHandler) SetupHandlers(server *echo.Echo) {
	server.GET("/api/post/:id/details", PostHandler.GetPost)
	server.POST("/api/post/:id/details", PostHandler.UpdatePost)
}
