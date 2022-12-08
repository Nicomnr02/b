package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/entity"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type ModifyWeb interface {
	AddTask(w http.ResponseWriter, r *http.Request)
	AddTaskProcess(w http.ResponseWriter, r *http.Request)
	AddCategory(w http.ResponseWriter, r *http.Request)
	AddCategoryProcess(w http.ResponseWriter, r *http.Request)

	UpdateTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskProcess(w http.ResponseWriter, r *http.Request)

	DeleteTask(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
}

type modifyWeb struct {
	taskClient     client.TaskClient
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewModifyWeb(tC client.TaskClient, cC client.CategoryClient, embed embed.FS) *modifyWeb {
	return &modifyWeb{tC, cC, embed}
}

func (a *modifyWeb) AddTask(w http.ResponseWriter, r *http.Request) {
	catId := r.URL.Query().Get("category")

	if r.Method == http.MethodGet {
		gHeader := path.Join("views", "general", "header.html")
		addCBody := path.Join("views", "main", "add-task.html")
		var templ = template.Must(template.ParseFS(a.embed, addCBody, gHeader))
		var err = templ.Execute(w, catId)
		if err != nil {
			panic(err) //tes
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("Method not allowed"))
		return
	}

	// ignore this
	_ = catId
	// TODO: answer here
}

func (a *modifyWeb) AddTaskProcess(w http.ResponseWriter, r *http.Request) {
	userId := fmt.Sprintf("%s", r.Context().Value("id"))

	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.URL.Query().Get("category")

	respCode, err := a.taskClient.CreateTask(title, description, category, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCode == 201 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	// ignore this
	_ = respCode

	// TODO: answer here
}

func (a *modifyWeb) AddCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: done
	if r.Method == http.MethodGet {
		gHeader := path.Join("views", "general", "header.html")
		addCBody := path.Join("views", "main", "add-category.html")
		var templ = template.Must(template.ParseFS(a.embed, addCBody, gHeader))
		var err = templ.Execute(w, nil)
		if err != nil {
			panic(err) //tes
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("Method not allowed"))
		return
	}
}

func (a *modifyWeb) AddCategoryProcess(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")
	category := r.FormValue("type")
	fmt.Println("add category process id user == ", userId)
	prID := userId.(string)

	respCode, err := a.categoryClient.AddCategories(category, prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("add category process id user string == ", prID)

	// ignore this
	_ = respCode
	//

	if respCode == 201 {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/category/add", http.StatusSeeOther)
	}
	// TODO: done
}

func (a *modifyWeb) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")

	task, err := a.taskClient.GetTaskById(taskId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ignore this
	_ = task
	//

	if r.Method == http.MethodGet {
		gHeader := path.Join("views", "general", "header.html")
		addCBody := path.Join("views", "main", "update-task.html")
		var templ = template.Must(template.ParseFS(a.embed, addCBody, gHeader))
		var err = templ.Execute(w, task)
		if err != nil {
			panic(err) //tes
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("Method not allowed"))
		return
	}

	// TODO: answer here
}

func (a *modifyWeb) UpdateTaskProcess(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")
	categoryId := r.URL.Query().Get("category_id")

	if categoryId == "" {
		title := r.FormValue("title")
		description := r.FormValue("description")

		respCode, err := a.taskClient.UpdateTask(taskId, title, description, r.Context().Value("id").(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if respCode == 200 {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/task/update?task_id="+taskId, http.StatusSeeOther)
		}
	} else {
		_, err := a.taskClient.UpdateCategoryTask(taskId, categoryId, r.Context().Value("id").(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

}

func (a *modifyWeb) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")

	respCode, err := a.taskClient.DeleteTask(taskId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if respCode == http.StatusOK {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	// TODO: answer here
}

func (a *modifyWeb) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")

	_, err := a.categoryClient.DeleteCategory(categoryId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: answer here
}
