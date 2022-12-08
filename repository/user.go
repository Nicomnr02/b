package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	userById := entity.User{}
	if err := r.db.Table("users").Select("*").Where("id = ?", id).Find(&userById).Error; err != nil {
		return entity.User{}, err
	}
	return userById, nil // TODO: get user by id
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	userByEmail := entity.User{}
	if err := r.db.Table("users").Select("*").Where("email = ?", email).Find(&userByEmail).Error; err != nil {
		return entity.User{}, err
	}
	return userByEmail, nil // TODO:get user by email
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil // TODO: CREATE USER
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	// TODO:UPDATE ACTION
	if err := r.db.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil

}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	if err := r.db.Table("users").Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil // TODO:DELETE ACTION
}
