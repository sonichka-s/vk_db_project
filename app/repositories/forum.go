package repositories

import (
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"vk_db_project/app/models"
)

type IForumRepository interface {
	CreateNewForum(models.Forum) (models.Forum, error)
	GetForum(string) (models.Forum, error)
	CreateThread(models.Thread) (models.Thread, error)
	GetThreads(models.Forum, int, string, bool) ([]models.Thread, error)
	GetForumUsers(string, int, string, bool) ([]models.UserModel, error)
}

type ForumRepoImpl struct {
	database *pgx.ConnPool
}

func NewForumRepoImpl(db *pgx.ConnPool) ForumRepoImpl {
	return ForumRepoImpl{database: db}
}

func (Forum ForumRepoImpl) CreateNewForum(forum models.Forum) (models.Forum, error) {
	userID := 0
	row := Forum.database.QueryRow("SELECT u_id , nickname FROM users WHERE nickname = $1", forum.User)

	err := row.Scan(&userID, &forum.User)
	if err != nil {
		return forum, err
	}

	_, err = Forum.database.Exec("INSERT INTO forums (slug , title, u_nickname) VALUES($1 , $2 , $3)", forum.Slug, forum.Title, forum.User)
	if err != nil {
		row := Forum.database.QueryRow("SELECT u_nickname , title , slug FROM forums WHERE slug = $1;", forum.Slug)
		row.Scan(&forum.User, &forum.Title, &forum.Slug)
		return forum, err
	}

	return forum, nil
}

func (Forum ForumRepoImpl) GetForum(slug string) (models.Forum, error) {

	forumData := new(models.Forum)
	row := Forum.database.QueryRow("SELECT slug , title, u_nickname , message_counter , thread_counter FROM forums WHERE slug = $1", slug)

	err := row.Scan(&forumData.Slug, &forumData.Title, &forumData.User, &forumData.Posts, &forumData.Threads)
	if err != nil {
		return *forumData, err
	}

	return *forumData, nil
}

func (Forum ForumRepoImpl) CreateThread(thread models.Thread) (models.Thread, error) {

	tx, err := Forum.database.Begin()

	if err != nil {
		return thread, err
	}

	userId := int64(0)
	var timer time.Time
	insertValues := make([]interface{}, 0)
	valuesCounter := 4
	valuesQuery := " VALUES($1 ,$2, $3, $4,"
	insertQuery := "INSERT INTO threads "
	insertColumns := "(message , title , u_nickname , f_slug ,"
	returningQuery := " RETURNING date , t_id"

	tx.Prepare("get-author", "SELECT u_id , nickname FROM users WHERE nickname = $1")

	row := tx.QueryRow("get-author", thread.Author)
	err = row.Scan(&userId, &thread.Author)
	if err != nil {
		tx.Rollback()
		return thread, err
	}

	tx.Prepare("get-forum", "SELECT slug FROM forums WHERE slug = $1")

	row = tx.QueryRow("get-forum", thread.Forum)
	err = row.Scan(&thread.Forum)
	if err != nil {
		tx.Rollback()
		return thread, err
	}

	insertValues = append(insertValues, thread.Message, thread.Title, thread.Author, thread.Forum)

	if thread.Slug != "" {
		insertColumns += " slug,"
		insertValues = append(insertValues, thread.Slug)
		valuesCounter++
		valuesQuery += " $" + strconv.Itoa(valuesCounter) + ","
	}

	if thread.Created.String() != "" {
		insertColumns += " date,"
		insertValues = append(insertValues, thread.Created)
		valuesCounter++
		valuesQuery += " $" + strconv.Itoa(valuesCounter) + ","
	}

	insertColumns = insertColumns[:len(insertColumns)-1] + ")"
	valuesQuery = valuesQuery[:len(valuesQuery)-1] + ")"

	err = tx.QueryRow(insertQuery+insertColumns+valuesQuery+returningQuery, insertValues...).Scan(&timer, &thread.Id)

	if timer.String() != "" {
		timer.Format(time.RFC3339)

		thread.Created = timer
	}

	if err != nil {
		tx.Rollback()
		row = Forum.database.QueryRow("SELECT u_nickname , date ,f_slug , t_id , message , slug , title , votes FROM threads WHERE slug = $1", thread.Slug)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		return thread, errors.New("thread already exist")
	}

	_, err = tx.Exec("INSERT INTO forumUsers (f_slug,u_nickname) VALUES ($1,$2) ON CONFLICT (f_slug,u_nickname) DO NOTHING", thread.Forum, thread.Author)
	_, err = tx.Exec("UPDATE forums SET thread_counter = thread_counter +1 WHERE slug = $1", thread.Forum)

	tx.Commit()
	return thread, nil
}

func (Forum ForumRepoImpl) GetThreads(forum models.Forum, limit int, since string, sort bool) ([]models.Thread, error) {

	tx, err := Forum.database.Begin()

	if err != nil {
		return nil , err
	}

	defer func () {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	} ()

	orderStatus := "DESC"
	sorter := "<"

	if !sort {
		sorter = ">"
		orderStatus = "ASC"
	}

	var rowThreads *pgx.Rows
	selectRow := "SELECT t_id , date , message , title , votes , slug , f_slug , u_nickname FROM threads T "
	if since != "" {
		sinceStatus := "WHERE f_slug = $3 AND date" + sorter + "=$2" + " "
		rowThreads, err = tx.Query(selectRow+sinceStatus+" ORDER BY date "+orderStatus+" LIMIT $1", limit, since, forum.Slug)
	} else {
		rowThreads, err = tx.Query(selectRow+"WHERE f_slug = $2 "+"ORDER BY date "+orderStatus+" LIMIT $1", limit, forum.Slug)
	}

	if err != nil {
		return nil, err
	}

	threads := make([]models.Thread, 0)

	if rowThreads != nil {

		for rowThreads.Next() {
			thread := new(models.Thread)
			var threadSlug *string
			err = rowThreads.Scan(&thread.Id, &thread.Created, &thread.Message, &thread.Title, &thread.Votes, &threadSlug, &thread.Forum, &thread.Author)

			if threadSlug != nil {
				thread.Slug = *threadSlug
			}

			if err != nil {
				return nil, err
			}

			threads = append(threads, *thread)
		}

		rowThreads.Close()
	}

	if len(threads) == 0 {
		tx.Prepare("get-slug","SELECT slug FROM forums WHERE slug = $1")
		row := tx.QueryRow("get-slug", forum.Slug)
		err = row.Scan(&forum.Slug)
	}

	return threads, err
}

func (Forum ForumRepoImpl) GetForumUsers(slug string, limit int, since string, desc bool) ([]models.UserModel, error) {

	var err error
	var row *pgx.Rows

	order := "DESC"
	ranger := "<"
	users := make([]models.UserModel, 0)

	if !desc {
		order = "ASC"
		ranger = ">"
	}

	selectRow := "SELECT U.nickname , U.fullname, U.email , U.about FROM  Users U WHERE U.nickname IN (SELECT FU.u_nickname FROM forumUsers FU WHERE FU.f_slug = $1 "
	selectValues := make([]interface{}, 0)
	if since != "" {
		if limit == 0 {
			selectRow += "AND FU.u_nickname " + ranger + " $2) ORDER BY U.nickname " + order
			selectValues = append(selectValues, slug, since)
		} else {
			selectRow += " AND FU.u_nickname " + ranger + " $3) ORDER BY U.nickname " + order + " LIMIT $2"
			selectValues = append(selectValues, slug, limit, since)
		}
	} else {
		if limit == 0 {
			selectRow += ") ORDER BY U.nickname " + order
			selectValues = append(selectValues, slug)
		} else {
			selectRow += ") ORDER BY U.nickname " + order + " LIMIT $2"
			selectValues = append(selectValues, slug, limit)
		}
	}

	row, err = Forum.database.Query(selectRow, selectValues...)

	if err != nil {
		return nil, err
	}

	if row != nil {
		for row.Next() {
			user := new(models.UserModel)
			err = row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)

			if err != nil {
				return nil, err
			}

			users = append(users, *user)
		}

		row.Close()
	}

	if len(users) == 0 {
		frow := Forum.database.QueryRow("SELECT slug FROM forums WHERE slug = $1", slug)
		err = frow.Scan(&slug)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
