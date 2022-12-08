package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	categories := []entity.Category{}

	rows, err := r.db.Table("categories").Select("categories.id,categories.type").Joins("inner join users on categories.user_id = users.id").Where("categories.user_id = ?", id).Rows()
	if err != nil {
		return []entity.Category{}, err
	}
	defer rows.Close()

	for rows.Next() {
		r.db.ScanRows(rows, &categories)

	}

	for _, c := range categories {
		fmt.Println("semua cate ==  ", c)
	}

	return categories, nil // TODO: done
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {

	if err := r.db.Save(&category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil // TODO: done
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {

	if err := r.db.Save(&categories).Error; err != nil {
		return err
	}
	return nil
	//todo: done
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	categoryByID := entity.Category{}
	if err := r.db.First(categoryByID, "id = ?", id).Error; err != nil {
		return entity.Category{}, err
	}

	return categoryByID, nil // TODO: done
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	if err := r.db.Model(&category).Where("id = ?", category.ID).Updates(&category).Error; err != nil {
		return err
	}
	return nil // TODO: done
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Category{}).Error; err != nil {
		return err
	}
	return nil // TODO: done
}
