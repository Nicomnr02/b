package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"fmt"

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
	rows, err := r.db.Table("tasks").Select("*").Joins("INNER JOIN users on tasks.user_id = users.id").Rows()
	if err != nil {
		return []entity.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		r.db.ScanRows(rows, &tasks)

	}

	fmt.Println(tasks, " tasks  -- asdfasfsadf")
	return tasks, nil

	// TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	if err := r.db.Save(&task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil // TODO: replace this
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	taskByID := entity.Task{ID: id}
	if err := r.db.Where("id = ?", taskByID.ID).First(&taskByID).Error; err != nil {
		return entity.Task{}, err
	}

	return taskByID, nil // TODO: done
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	taskByCatID := []entity.Task{}
	if err := r.db.Find(&taskByCatID, "category_id = ?", catId).Error; err != nil {
		return []entity.Task{}, err
	}
	return taskByCatID, nil // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {

	if err := r.db.Model(task).Where("ID = ?", task.ID).Updates(task).Error; err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Task{}).Error; err != nil {
		return err
	}
	return nil // TODO: replace this
}
