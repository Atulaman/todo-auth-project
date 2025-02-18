package handler

import (
	"database/sql"

	"net/http"
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
	err := utils.DecodeJson(r, &newTask)
	if err != nil || newTask.Desc == "" {
		utils.ResponseError(w, "Invalid request", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	username, err := dbhelper.GetUser(r)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseError(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	newTask.Id, err = dbhelper.GetTaskId(username)
	if err != nil {
		utils.ResponseError(w, "Error while generating id", http.StatusInternalServerError)
		log.Logging(err, "Error while generating task id", 500, "error", r)
		return
	}
	err = dbhelper.CreateTask(username, newTask.Desc, newTask.Id)
	if err != nil {
		utils.ResponseError(w, "Error while adding task", http.StatusInternalServerError)
		log.Logging(err, "Error while adding task", 500, "error", r)
		return
	}
	utils.ResponseJson(w, http.StatusCreated, map[string]interface{}{"message": "Task added successfully", "task": newTask})
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
	tasks, err := dbhelper.GetTask(r)
	if err != nil {
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error while fetching tasks", 500, "error", r)
		return
	}
	if len(tasks) == 0 {
		utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"message": "No tasks found", "count": 0, "tasks": []task{}})
		log.Logging(err, "No tasks found", 200, "info", r)
		return
	}
	utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"tasks": tasks, "message": "success", "count": len(tasks)})
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
	err := utils.DecodeJson(r, &newTask)
	if err != nil || newTask.Id <= 0 || newTask.Desc == "" {
		utils.ResponseError(w, "Invalid Id or Description", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	// cookie, _ := utils.GetSessionID(r)
	// var username string
	// query := `SELECT username FROM session WHERE session_id=$1`
	// err = database.TODO.Get(&username, query, cookie)
	username, err := dbhelper.GetUser(r)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseError(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	//result, err := database.TODO.Exec(`UPDATE "Tasks" SET description=$2 WHERE id=$1 AND username=$3 AND archive = false`, newTask.Id, newTask.Desc, username)
	//fmt.Printf("type of result=%T", result)
	result, err := dbhelper.UpdateTask(newTask.Id, newTask.Desc, username)
	if err != nil {
		utils.ResponseError(w, "Error while updating task", http.StatusInternalServerError)
		log.Logging(err, "Error while updating task", 500, "error", r)
		return
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.ResponseError(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	if RowsAffected == 0 {
		utils.ResponseJson(w, http.StatusNotFound, map[string]interface{}{"message": "Task not found"})
		log.Logging(err, "Task not found", 404, "warning", r)
		return
	}
	utils.ResponseJson(w, http.StatusOK, map[string]interface{}{"message": "Task Updated successfully", "task": newTask})
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
	err := utils.DecodeJson(r, &newTask)
	if err != nil || newTask.Id <= 0 {
		utils.ResponseError(w, "Invalid request", http.StatusBadRequest)
		log.Logging(err, "Invalid request", 400, "warning", r)
		return
	}
	// 	cookie, _ := utils.GetSessionID(r)
	// 	var username string
	// 	query := `SELECT username
	// FROM session
	// WHERE session_id = $1`
	// 	err = database.TODO.Get(&username, query, cookie)
	username, err := dbhelper.GetUser(r)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseError(w, "Unauthorized user", http.StatusUnauthorized)
			log.Logging(err, "Unauthorized user", 401, "warning", r)
			return
		}
		utils.ResponseError(w, err.Error(), http.StatusInternalServerError)
		log.Logging(err, "Error checking user", 500, "error", r)
		return
	}
	//result, err := database.TODO.Exec(`UPDATE "Tasks" SET archive = true WHERE id=$1 AND username=$2 AND archive = false`, newTask.Id, username)
	result, err := dbhelper.DeleteTask(newTask.Id, username)
	if err != nil {
		utils.ResponseError(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.ResponseError(w, "Database error", http.StatusInternalServerError)
		log.Logging(err, "Database error", 500, "error", r)
		return
	}
	if RowsAffected == 0 {
		utils.ResponseJson(w, http.StatusNotFound, map[string]string{"message": "Task not found!"})
		log.Logging(err, "Task not found", 404, "warning", r)
		return
	}
	utils.ResponseJson(w, http.StatusOK, map[string]string{"message": "Task has been deleted successfully!"})
	log.Logging(err, "Task deleted successfully", 200, "info", r)
}
