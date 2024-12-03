package handler

import (
	"database/sql"
	"encoding/json"

	//"log"
	"net/http"
	"todo-auth/database"
	dbhelper "todo-auth/database/db-helper"
	log "todo-auth/logging"
	"todo-auth/utils"

	_ "github.com/lib/pq"
)

type task struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

// Add godoc
// @Summary Add a new task
// @Description Add a new task for the logged-in user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body task true "task to add"
// @Success 200 {object} map[string]interface{} "Task added successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Error adding task"
// @Router /tasks [post]
func Add(w http.ResponseWriter, r *http.Request) {
	var newTask task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil || newTask.Desc == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	username, err := dbhelper.GetUser(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	newTask.Id, err = dbhelper.GetTaskId(username)
	if err != nil {
		http.Error(w, "Error while generating id", http.StatusInternalServerError)
		log.Logging(err, "Error while generating task id", 500, "error", r)
		return
	}
	err = dbhelper.CreateTask(username, newTask.Desc, newTask.Id)
	if err != nil {
		http.Error(w, "Error while adding task", http.StatusInternalServerError)
		log.Logging(err, "Error while adding task", 500, "error", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Task added successfully", "task": newTask})
	log.Logging(err, "Task added successfully", 201, "info", r)
}

// List godoc
// @Summary List all tasks
// @Description Get all tasks for the logged-in user
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} []task "tasks fetched successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Error fetching tasks"
// @Router /tasks [get]
func List(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.GetSessionID(r)
	query := `SELECT t1.id, t1.description
					FROM "Tasks" t1
							JOIN session a ON t1.username = a.username
					where a.session_id = $1 AND t1.archive = false
					ORDER BY t1.id ASC `
	tasks := []struct {
		Id   int    `db:"id"`
		Desc string `db:"description"`
	}{}
	err := database.TODO.Select(&tasks, query, cookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error while fetching tasks", 500, "error", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(tasks) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "No tasks found", "count": 0, "tasks": []task{}})
		log.Logging(err, "No tasks found", 200, "info", r)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks, "message": "success", "count": len(tasks)})
	log.Logging(err, "Tasks fetched successfully", 200, "info", r)
}

// Update godoc
// @Summary Update a task
// @Description Update the description of an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body task true "task to update"
// @Success 200 {object} map[string]interface{} "task updated successfully"
// @Failure 400 {object} map[string]string "Invalid task ID or description"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "task not found"
// @Failure 500 {object} map[string]string "Error updating task"
// @Router /tasks [put]
func Update(w http.ResponseWriter, r *http.Request) {
	var newTask task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil || newTask.Id <= 0 || newTask.Desc == "" {
		http.Error(w, "Invalid Id or Description", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	cookie, _ := utils.GetSessionID(r)
	var username string
	query := `SELECT username FROM session WHERE session_id=$1`
	err = database.TODO.Get(&username, query, cookie)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	result, err := database.TODO.Exec(`UPDATE "Tasks" SET description=$2 WHERE id=$1 AND username=$3 AND archive = false`, newTask.Id, newTask.Desc, username)
	if err != nil {
		http.Error(w, "Error while updating task", http.StatusInternalServerError)
		log.Logging(err, "Error while updating task", 500, "error", r)
		return
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	if RowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Logging(err, "Task not found", 404, "warning", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Task Updated successfully", "task": newTask})
	log.Logging(err, "Task updated successfully", 200, "info", r)
}

// Delete godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body task true "task to delete"
// @Success 200 {object} map[string]string "task deleted successfully"
// @Failure 400 {object} map[string]string "Invalid task ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "task not found"
// @Failure 500 {object} map[string]string "Error deleting task"
// @Router /tasks [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	var newTask task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil || newTask.Id <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	cookie, _ := utils.GetSessionID(r)
	var username string
	query := `SELECT username
FROM session
WHERE session_id = $1`
	err = database.TODO.Get(&username, query, cookie)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	// 	query = `DELETE
	// FROM "Tasks"
	// WHERE id = $1
	//   AND username = $2`
	// 	result, err := database.TODO.Exec(query, newTask.Id, username)
	result, err := database.TODO.Exec(`UPDATE "Tasks" SET archive = true WHERE id=$1 AND username=$2 AND archive = false`, newTask.Id, username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	if RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Task not found!"})
		log.Logging(err, "Task not found", 404, "warning", r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task has been deleted successfully!"})
	log.Logging(err, "Task deleted successfully", 200, "info", r)
}
