package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type StudentRepository struct {
	db *sqlx.DB
}

func NewStudentRepository(db database.Db) repository.StudentRepository {
	return &StudentRepository{db: db}
}

func (repo StudentRepository) GetAllStudentQuery(c context.Context) []entity.Student {
	return []entity.Student{
		{Name: "John"},
		{Name: "Khoa"},
		{Name: "Lindan"},
	}
}
