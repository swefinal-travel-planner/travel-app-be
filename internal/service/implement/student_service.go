package serviceimplement

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type StudentService struct {
	studentRepository repository.StudentRepository
}

func NewStudentService(studentRepository repository.StudentRepository) service.StudentService {
	return &StudentService{studentRepository: studentRepository}
}

func (service StudentService) GetAllStudent(ctx context.Context) []model.Student {
	studentsFromRepo := service.studentRepository.GetAllStudentQuery(ctx)

	students := make([]model.Student, len(studentsFromRepo))
	for i, s := range studentsFromRepo {
		students[i] = model.Student{
			Name: s.Name,
		}
	}
	return students
}
