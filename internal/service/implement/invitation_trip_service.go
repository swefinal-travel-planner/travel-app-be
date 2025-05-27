package serviceimplement

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type InvitationTripService struct {
	invitationTripRepo  repository.InvitationTripRepository
	tripRepo            repository.TripRepository
	tripMemberRepo      repository.TripMemberRepository
	unitOfWork          repository.UnitOfWork
	notificationService service.NotificationService
}

func NewInvitationTripService(
	invitationTripRepo repository.InvitationTripRepository,
	tripRepo repository.TripRepository,
	tripMemberRepo repository.TripMemberRepository,
	unitOfWork repository.UnitOfWork,
	notificationService service.NotificationService,
) service.InvitationTripService {
	return &InvitationTripService{
		invitationTripRepo:  invitationTripRepo,
		tripRepo:            tripRepo,
		tripMemberRepo:      tripMemberRepo,
		unitOfWork:          unitOfWork,
		notificationService: notificationService,
	}
}

func (s *InvitationTripService) SendInvitation(ctx *gin.Context, invitation model.InvitationTripRequest, senderId int64) string {
	// Start transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation Begin transaction error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer s.unitOfWork.Rollback(tx)

	// Check if sender is a member of the trip
	isMember, err := s.tripMemberRepo.IsUserInTripQuery(ctx, invitation.TripID, senderId, tx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation IsUserInTripQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Lock the trip for update
	trip, err := s.tripRepo.SelectForUpdateById(ctx, invitation.TripID, tx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation SelectForUpdateById error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if trip == nil {
		return error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// Check if user is already invited
	existingInvitation, err := s.invitationTripRepo.GetOneByReceiverIdAndTripIDQuery(ctx, invitation.ReceiverID, invitation.TripID, tx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation GetByTripIdQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	for existingInvitation != nil {
		return error_utils.ErrorCode.TRIP_INVITATION_ALREADY_EXISTS
	}

	// Create invitation
	newInvitation := &entity.InvitationTrip{
		TripID:     invitation.TripID,
		SenderID:   senderId,
		ReceiverID: invitation.ReceiverID,
		Status:     "pending",
	}

	err = s.invitationTripRepo.CreateCommand(ctx, newInvitation, tx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation CreateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InvitationTripService.SendInvitation Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	errCode := s.notificationService.SaveAndSendNotification(ctx, model.SaveNotificationRequest{
		Type:                entity.NotificationType.FriendRequestReceived,
		ReceiverUserID:      invitation.ReceiverID,
		TriggerEntityType:   entity.NotificationTriggerType.User,
		TriggerEntityID:     &senderId,
		ReferenceEntityType: entity.NotificationReferenceType.FriendInvitation,
		ReferenceEntityID:   &newInvitation.ID,
	})

	return errCode
}

func (s *InvitationTripService) GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationTripReceivedResponse, string) {
	invitations, err := s.invitationTripRepo.GetByReceiverIDQuery(ctx, userId, nil)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return make([]model.InvitationTripReceivedResponse, 0), ""
		}
		log.Error("InvitationTripService.GetAllReceivedInvitations GetByReceiverIdQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	response := make([]model.InvitationTripReceivedResponse, len(invitations))
	for i, inv := range invitations {
		response[i] = model.InvitationTripReceivedResponse{
			ID:         inv.ID,
			TripID:     inv.TripID,
			SenderID:   inv.SenderID,
			ReceiverID: inv.ReceiverID,
			Status:     inv.Status,
			CreatedAt:  inv.CreatedAt,
			UpdatedAt:  inv.UpdatedAt,
		}
	}

	return response, ""
}

func (s *InvitationTripService) GetAllSentInvitations(ctx *gin.Context, userId int64) ([]model.InvitationTripSentResponse, string) {
	invitations, err := s.invitationTripRepo.GetBySenderIDQuery(ctx, userId, nil)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return make([]model.InvitationTripSentResponse, 0), ""
		}
		log.Error("InvitationTripService.GetAllSentInvitations GetBySenderIdQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}

	response := make([]model.InvitationTripSentResponse, len(invitations))
	for i, inv := range invitations {
		response[i] = model.InvitationTripSentResponse{
			ID:         inv.ID,
			TripID:     inv.TripID,
			SenderID:   inv.SenderID,
			ReceiverID: inv.ReceiverID,
			Status:     inv.Status,
			CreatedAt:  inv.CreatedAt,
			UpdatedAt:  inv.UpdatedAt,
		}
	}

	return response, ""
}

func (s *InvitationTripService) AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	// Start transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InvitationTripService.AcceptInvitation Begin transaction error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer s.unitOfWork.Rollback(tx)

	// Get invitation
	invitation, err := s.invitationTripRepo.GetOneByIDQuery(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.AcceptInvitation GetByIdQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if invitation == nil {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if user is the receiver
	if invitation.ReceiverID != userId {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if invitation is pending
	if invitation.Status != "pending" {
		return error_utils.ErrorCode.FORBIDDEN
	}

	err = s.invitationTripRepo.DeleteByIDCommand(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.AcceptInvitation UpdateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Add user to trip members
	tripMember := &entity.TripMember{
		TripID: invitation.TripID,
		UserID: userId,
	}
	err = s.tripMemberRepo.CreateCommand(ctx, tripMember, tx)
	if err != nil {
		log.Error("InvitationTripService.AcceptInvitation CreateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InvitationTripService.AcceptInvitation Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (s *InvitationTripService) DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	// Start transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InvitationTripService.DenyInvitation Begin transaction error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer s.unitOfWork.Rollback(tx)

	// Get invitation
	invitation, err := s.invitationTripRepo.GetOneByIDQuery(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.DenyInvitation GetByIdQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if invitation == nil {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if user is the receiver
	if invitation.ReceiverID != userId {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if invitation is pending
	if invitation.Status != "pending" {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Update invitation status
	err = s.invitationTripRepo.DeleteByIDCommand(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.DenyInvitation UpdateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InvitationTripService.DenyInvitation Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (s *InvitationTripService) WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	// Start transaction
	tx, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("InvitationTripService.WithdrawInvitation Begin transaction error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer s.unitOfWork.Rollback(tx)

	// Get invitation
	invitation, err := s.invitationTripRepo.GetOneByIDQuery(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.WithdrawInvitation GetByIdQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if invitation == nil {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if user is the sender
	if invitation.SenderID != userId {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Check if invitation is pending
	if invitation.Status != "pending" {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Update invitation status
	err = s.invitationTripRepo.DeleteByIDCommand(ctx, invitationId, tx)
	if err != nil {
		log.Error("InvitationTripService.WithdrawInvitation UpdateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Commit transaction
	err = s.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("InvitationTripService.WithdrawInvitation Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}
