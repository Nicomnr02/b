package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	currentUserID := fmt.Sprintf("%s", r.Context().Value("id"))
	currentUserIDint, err := strconv.Atoi(currentUserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		// panic(err)
		return
	}

	if taskID == "" {
		if tasks1, err := t.taskService.GetTasks(r.Context(), currentUserIDint); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&tasks1)
		}
	} else {

		taskIDint, err := strconv.Atoi(taskID)
		if err != nil {
			panic("error here")
		}
		if tasks, err := t.taskService.GetTaskByID(r.Context(), taskIDint); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			fmt.Print("error here")
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&tasks)
		}
	}

	// TODO: done
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	userID := fmt.Sprintf("%s", r.Context().Value("id"))
	intID, _ := strconv.Atoi(userID)
	if newTask, err := t.taskService.StoreTask(r.Context(), &entity.Task{ID: task.ID, Title: task.Title, Description: task.Description, CategoryID: task.CategoryID, UserID: intID}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID, "task_id": newTask.ID, "message": "success create new task"})
	}

	// TODO: done
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	taskIDint, _ := strconv.Atoi(taskID)
	userID := r.Context().Value("id")

	if err := t.taskService.DeleteTask(r.Context(), taskIDint); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		floatID := float64(taskIDint)
		json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID, "task_id": floatID, "message": "success delete task"})
	}

	// TODO: done
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	currentUserID := fmt.Sprintf("%s", r.Context().Value("id"))
	if currentUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	if newTask, err := t.taskService.UpdateTask(r.Context(), &entity.Task{ID: task.ID, Title: task.Title, Description: task.Description, CategoryID: task.CategoryID}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"user_id": currentUserID, "task_id": newTask.ID, "message": "success update task"})
	}

	// TODO: done
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
