package serviceimplement

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type TripMemberService struct {
	tripMemberRepo repository.TripMemberRepository
}

func NewTripMemberService(tripMemberRepo repository.TripMemberRepository) service.TripMemberService {
	return &TripMemberService{
		tripMemberRepo: tripMemberRepo,
	}
}

func (s *TripMemberService) GetTripMembersIfUserInTrip(ctx context.Context, tripID int64, userID int64) ([]model.TripMemberResponse, string) {
	isMember, err := s.tripMemberRepo.IsUserInTripQuery(ctx, tripID, userID, nil)
	if err != nil {
		log.Error("TripMemberService.GetTripMembersIfUserInTrip IsUserInTripQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	if !isMember {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

	members, err := s.tripMemberRepo.GetTripMembersQuery(ctx, tripID, nil)
	if err != nil {
		log.Error("TripMemberService.GetTripMembersIfUserInTrip GetTripMembersQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	responses := make([]model.TripMemberResponse, len(members))
	for i, member := range members {
		responses[i] = model.TripMemberResponse{
			ID:       member.ID,
			TripID:   member.TripID,
			UserID:   member.UserID,
			Role:     member.Role,
			Name:     member.Name,
			PhotoURL: member.PhotoURL,
		}
	}

	return responses, ""
}

func (s *TripMemberService) DeleteMemberFromTrip(ctx context.Context, tripID int64, memberID int64, deleterID int64) string {
	if deleterID == memberID {
		return error_utils.ErrorCode.FORBIDDEN
	}
	// Check deleter is admin
	isAdmin, err := s.tripMemberRepo.IsUserTripAdminQuery(ctx, tripID, deleterID, nil)
	if err != nil {
		log.Error("TripMemberService.DeleteMemberFromTrip IsUserTripAdminQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdmin {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check member is in trip
	isMember, err := s.tripMemberRepo.IsUserInTripQuery(ctx, tripID, memberID, nil)
	if err != nil {
		log.Error("TripMemberService.DeleteMemberFromTrip IsUserInTripQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return error_utils.ErrorCode.FORBIDDEN
	}

	err = s.tripMemberRepo.DeleteMemberCommand(ctx, tripID, memberID, nil)
	if err != nil {
		log.Error("TripMemberService.DeleteMemberFromTrip DeleteMemberCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	return ""
}
