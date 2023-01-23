package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	"vk_db_project/app/handlers"
	repos "vk_db_project/app/repositories"
	usecases "vk_db_project/app/uscases"

	"github.com/jackc/pgx"
)

const (
	usernameDB = "docker"
	passwordDB = "docker"
	nameDB     = "docker"
)

type RequestHandler struct {
	userHandler   handlers.UserHandler
	forumHandler  handlers.ForumHandler
	threadHandler handlers.ThreadHandler
	postHandler   handlers.PostHandler
}

func StartServer(db *pgx.ConnPool) *RequestHandler {

	postDB := repos.NewPostRepoImpl(db)
	postUse := usecases.NewPostUsecaseImpl(postDB)
	postH := handlers.NewPostHandler(postUse)

	threadDB := repos.NewThreadRepoImpl(db)
	threadUse := usecases.NewThreadsUsecaseImpl(threadDB)
	threadH := handlers.NewThreadHandler(threadUse)

	forumDB := repos.NewForumRepoImpl(db)
	forumUse := usecases.NewForumUsecaseImpl(forumDB)
	forumH := handlers.NewForumHandler(forumUse)

	userDB := repos.NewUserRepoImpl(db)
	userUse := usecases.NewUserUsecaseImpl(userDB)
	userH := handlers.NewUserHandler(userUse)

	api := &RequestHandler{userHandler: userH, forumHandler: forumH, threadHandler: threadH, postHandler: postH}

	return api
}

func JSONMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		return next(c)
	}
}

func Logs(next echo.HandlerFunc) echo.HandlerFunc {

	return func(rwContext echo.Context) error {

		var err error
		if rwContext.Request().Method == "GET" {
			start := time.Now()
			err = next(rwContext)
			respTime := time.Since(start)
			if respTime.Milliseconds() >= 400 {
				fmt.Println("MICRO SEC:", respTime.Microseconds(), "\n PATH:", rwContext.Request().URL.Path, "\n METHOD:", rwContext.Request().Method)
				fmt.Println(rwContext.QueryParam("sort"))
			}

		} else {
			err = next(rwContext)
		}

		return err

	}
}

func main() {
	server := echo.New()
	connectString := "user=" + usernameDB + " password=" + passwordDB + " dbname=" + nameDB + " sslmode=disable"

	pgxConn, err := pgx.ParseConnectionString(connectString)
	pgxConn.PreferSimpleProtocol = false

	if err != nil {
		server.Logger.Fatal("PARSING CONFIG ERROR", err.Error())
	}

	config := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	connPool, err := pgx.NewConnPool(config)
	defer connPool.Close()
	if err != nil {
		server.Logger.Fatal("NO CONNECTION TO BD", err.Error())
	}
	fmt.Println(connPool.Stat())
	api := StartServer(connPool)
	api.userHandler.SetupHandlers(server)
	api.forumHandler.SetupHandlers(server)
	api.threadHandler.SetupHandlers(server)
	api.postHandler.SetupHandlers(server)

	server.Logger.Fatal(server.Start(":5000"))
}
