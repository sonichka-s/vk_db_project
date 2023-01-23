package repositories

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx"
	"vk_db_project/app/models"
)

type IUserRepo interface {
	CreateNewUser(models.UserModel) ([]models.UserModel, error)
	UpdateUserData(models.UserModel) (models.UserModel, error)
	GetUserData(string) (models.UserModel, error)
	Status() models.Status
	Clear()
}

type UserRepoImpl struct {
	database *pgx.ConnPool
}

func NewUserRepoImpl(db *pgx.ConnPool) UserRepoImpl {
	return UserRepoImpl{database: db}
}

func (User UserRepoImpl) CreateNewUser(userModel models.UserModel) ([]models.UserModel, error) {

	allData := make([]models.UserModel, 0)
	var err error

	_, err = User.database.Exec("INSERT INTO users (nickname , fullname , email , about) VALUES($1 , $2 , $3 ,$4)", userModel.Nickname, userModel.Fullname, userModel.Email, userModel.About)

	if err != nil {
		row, err := User.database.Query("SELECT nickname , fullname , email , about FROM users WHERE nickname = $1 OR email = $2", userModel.Nickname, userModel.Email)

		if row != nil {
			for row.Next() {

				if err == nil {
					err = errors.New("such user already exists")
				}

				existingUser := models.UserModel{
					Nickname: "",
					Fullname: "",
					Email:    "",
					About:    "",
				}

				row.Scan(&existingUser.Nickname, &existingUser.Fullname, &existingUser.Email, &existingUser.About)

				allData = append(allData, existingUser)
			}

			row.Close()
		}

		return allData, errors.New("such user already exists")
	}

	allData = append(allData, userModel)

	return allData, err
}

func (User UserRepoImpl) UpdateUserData(userModel models.UserModel) (models.UserModel, error) {

	id := 2
	values := make([]interface{}, 0)

	querySting := "UPDATE users SET"
	nickQuery := " WHERE nickname = $1 RETURNING u_id, nickname, fullname , email, about"
	reqQuery := ""

	values = append(values, userModel.Nickname)

	if userModel.Email != "" {
		values = append(values, userModel.Email)
		reqQuery += " " + "email = $" + strconv.Itoa(id) + ","
		id++
	}

	if userModel.Fullname != "" {
		values = append(values, userModel.Fullname)
		reqQuery += " " + "fullname = $" + strconv.Itoa(id) + ","
		id++
	}

	if userModel.About != "" {
		values = append(values, userModel.About)
		reqQuery += " " + "about = $" + strconv.Itoa(id) + ","
		id++
	}

	if len(reqQuery) > 1 {
		reqQuery = reqQuery[:len(reqQuery)-1]
	}

	var row *pgx.Row

	if len(values) == 1 {
		row = User.database.QueryRow("SELECT u_id, nickname, fullname , email, about FROM users WHERE nickname = $1", values[0])
	} else {
		row = User.database.QueryRow(querySting+reqQuery+nickQuery, values...)
	}

	userId := 0

	err := row.Scan(&userId, &userModel.Nickname, &userModel.Fullname, &userModel.Email, &userModel.About)

	return userModel, err

}

func (User UserRepoImpl) GetUserData(nickname string) (models.UserModel, error) {

	userData := models.UserModel{
		Nickname: "",
		Fullname: "",
		Email:    "",
		About:    "",
	}

	row := User.database.QueryRow("SELECT nickname , fullname , email, about FROM users WHERE nickname = $1", nickname)

	err := row.Scan(&userData.Nickname, &userData.Fullname, &userData.Email, &userData.About)

	return userData, err
}

func (User UserRepoImpl) Status() models.Status {

	statAnswer := new(models.Status)
	row := User.database.QueryRow("SELECT (SELECT COUNT(u_id) FROM users) as uc , (SELECT COUNT(f_id) FROM forums) AS fc , (SELECT COUNT(t_id) FROM threads) AS tc , (SELECT COUNT(m_id) FROM messages) AS mc")
	row.Scan(&statAnswer.User, &statAnswer.Forum, &statAnswer.Thread, &statAnswer.Post)

	return *statAnswer
}

func (User UserRepoImpl) Clear() {
	User.database.Exec("DELETE FROM users;")
	User.database.Exec("DELETE FROM forums;")
	User.database.Exec("DELETE FROM threads;")
	User.database.Exec("DELETE FROM messages;")
	User.database.Exec("DELETE FROM voteThreads;")
	User.database.Exec("DELETE FROM forumUsers;")
}