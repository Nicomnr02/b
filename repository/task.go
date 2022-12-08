package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	tasks := []entity.Task{}

	if err := r.db.Table("tasks").Joins("INNER JOIN users on tasks.user_id = users.id").Where("tasks.user_id = ?", id).Find(&tasks).Error; err != nil {
		return []entity.Task{}, err
	}

	return tasks, nil

	// TODO: done
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	if err := r.db.Save(&task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil // TODO: done
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	taskByID := entity.Task{}
	if err := r.db.First(&taskByID, id).Error; err != nil {
		return entity.Task{}, err
	}

	return taskByID, nil // TODO: done
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	taskByCatID := []entity.Task{}
	if err := r.db.Find(&taskByCatID, "category_id = ?", catId).Error; err != nil {
		return []entity.Task{}, err
	}
	return taskByCatID, nil // TODO: done
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {

	if err := r.db.Model(&task).Where("ID = ?", task.ID).Updates(&task).Error; err != nil {
		return err
	}
	return nil // TODO: done
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Task{}).Error; err != nil {
		return err
	}
	return nil // TODO: done
}
