package dbhelper

import (
	"net/http"
	"todo-auth/database"
	"todo-auth/utils"
)

func GetTaskId(username string) (int, error) {
	p, err := database.TODO.Query(`SELECT
	CASE
	WHEN (SELECT id FROM "Tasks" WHERE id=1 AND username=$1) IS NULL THEN 1
	ELSE
	(select coalesce(min(t1.id +1),1) from (SELECT id FROM "Tasks" WHERE username=$1) t1 left join (SELECT id FROM "Tasks" WHERE username=$1) t2 on t1.id +1 =t2.id  where t2.id is null)
	END
	`, username)
	if err != nil {
		return 0, err
	}
	p.Next()
	var id int
	if err := p.Scan(&id); err != nil {
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
	err = database.TODO.QueryRow(`SELECT username FROM session WHERE session_id=$1`, cookie).Scan(&username)
	return username, err
}

func CreateTask(username string, desc string, id int) error {
	_, err := database.TODO.Exec(`INSERT INTO "Tasks" (id,description,username) VALUES ($1,$2,$3)`, id, desc, username)
	return err
}
