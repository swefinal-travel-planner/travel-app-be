package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type StudentRepository interface {
	GetAllStudentQuery(c context.Context) []entity.Student
}
