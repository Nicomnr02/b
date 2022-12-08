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

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctxID := fmt.Sprintf("%s", r.Context().Value("id"))

	catID, err := strconv.Atoi(ctxID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	if categories, err := c.categoryService.GetCategories(r.Context(), catID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&categories)
	}
	// TODO: done
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	currentUserID := fmt.Sprintf("%s", r.Context().Value("id"))
	if currentUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	currUserIDint, _ := strconv.Atoi(currentUserID)
	if newCat, err := c.categoryService.StoreCategory(r.Context(), &entity.Category{Type: category.Type, UserID: currUserIDint}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	} else {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{"user_id": currentUserID, "category_id": newCat.ID, "message": "success create new category"})
	}

	// TODO: done
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	catID, _ := strconv.Atoi(categoryID)
	currentUserID := fmt.Sprintf("%s", r.Context().Value("id"))

	if err := c.categoryService.DeleteCategory(r.Context(), catID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": currentUserID, "category_id": catID, "message": "success delete category"})
	// TODO: done
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&categories)

}
