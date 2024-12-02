package dbhelper

import (
	"net/http"
	"todo-auth/database"
	"todo-auth/utils"
)

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
