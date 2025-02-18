package dbhelper

import (
	"database/sql"
	"net/http"
	"todo-auth/database"
	"todo-auth/utils"
)

// type TaskGet struct {
// 	Id   int    `db:"id"`
// 	Desc string `db:"description"`
// }

func GetTaskId(username string) (int, error) {
	query := `SELECT
	CASE
	WHEN (SELECT id FROM "Tasks" WHERE id=1 AND username=$1) IS NULL THEN 1
	ELSE
	(select coalesce(min(t1.id +1),1) from (SELECT id FROM "Tasks" WHERE username=$1) t1 left join (SELECT id FROM "Tasks" WHERE username=$1) t2 on t1.id +1 =t2.id  where t2.id is null)
	END`
	var id int
	err := database.TODO.Get(&id, query, username)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func GetUser(r *http.Request) (string, error) {
	cookie, err := utils.GetSessionID(r)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.ErrNoCookie
		}
		return "", err
	}
	var username string = ""
	query := `SELECT username FROM session WHERE session_id=$1`
	err = database.TODO.Get(&username, query, cookie)
	return username, err
}

func CreateTask(username string, desc string, id int) error {
	query := `INSERT INTO "Tasks" (id,description,username) VALUES ($1,$2,$3)`
	_, err := database.TODO.Exec(query, id, desc, username)
	return err
}

// list controller
func GetTask(r *http.Request) (tasks []struct {
	Id   int    `db:"id"`
	Desc string `db:"description"`
}, err error) {
	cookie, _ := utils.GetSessionID(r)
	query := `SELECT t1.id, t1.description
					FROM "Tasks" t1
							JOIN session a ON t1.username = a.username
					where a.session_id = $1 AND t1.archive = false
					ORDER BY t1.id ASC `

	err = database.TODO.Select(&tasks, query, cookie)
	return
}

// UPDATE
func UpdateTask(id int, desc string, username string) (result sql.Result, err error) {
	query := `UPDATE "Tasks" SET description=$2 WHERE id=$1 AND username=$3 AND archive = false`
	result, err = database.TODO.Exec(query, id, desc, username)
	//return result, err
	return
}

// DELETE
func DeleteTask(id int, username string) (result sql.Result, err error) {
	query := `UPDATE "Tasks" SET archive = true WHERE id=$1 AND username=$2 AND archive = false`
	result, err = database.TODO.Exec(query, id, username)
	return
}
