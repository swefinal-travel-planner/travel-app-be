package service

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type StudentService interface {
	GetAllStudent(ctx context.Context) []model.Student
}
