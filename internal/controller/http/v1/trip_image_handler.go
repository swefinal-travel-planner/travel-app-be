package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type TripImageHandler struct {
	tripImageService service.TripImageService
}

func NewTripImageHandler(tripImageService service.TripImageService) *TripImageHandler {
	return &TripImageHandler{
		tripImageService: tripImageService,
	}
}

// @Summary Create trip image
// @Description Add a new image to a trip
// @Tags TripImages
// @Accept json
// @Produce json
// @Param tripId path int true "Trip ID"
// @Param request body model.TripImageRequest true "Trip image payload"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /trips/{tripId}/images [post]
func (h *TripImageHandler) CreateTripImage(c *gin.Context) {
	userID := middleware.GetUserIdHelper(c)

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		c.JSON(statusCode, errResponse)
		return
	}

	var tripImageRequest model.TripImageRequest
	if err := validation.BindJsonAndValidate(c, &tripImageRequest); err != nil {
		return
	}

	errCode := h.tripImageService.CreateTripImage(c, userID, tripID, tripImageRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		c.JSON(statusCode, errResponse)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary Get trip images
// @Description Get all images for a trip with author information
// @Tags TripImages
// @Accept json
// @Produce json
// @Param tripId path int true "Trip ID"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.TripImageWithUserInfoResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /trips/{tripId}/images [get]
func (h *TripImageHandler) GetTripImages(c *gin.Context) {
	userID := middleware.GetUserIdHelper(c)

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		c.JSON(statusCode, errResponse)
		return
	}

	tripImagesWithUserInfo, errCode := h.tripImageService.GetTripImagesWithUserInfo(c, userID, tripID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		c.JSON(statusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&tripImagesWithUserInfo))
}

// @Summary Delete trip image
// @Description Delete a specific image from a trip (admin only)
// @Tags TripImages
// @Accept json
// @Produce json
// @Param tripId path int true "Trip ID"
// @Param imageId path int true "Image ID"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /trips/{tripId}/images/{imageId} [delete]
func (h *TripImageHandler) DeleteTripImage(c *gin.Context) {
	userID := middleware.GetUserIdHelper(c)

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		c.JSON(statusCode, errResponse)
		return
	}

	imageID, err := strconv.ParseInt(c.Param("imageId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "imageId")
		c.JSON(statusCode, errResponse)
		return
	}

	errCode := h.tripImageService.DeleteTripImage(c, userID, tripID, imageID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		c.JSON(statusCode, errResponse)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary Get all trip images by trip ID and trip item ID
// @Description Get all images for a trip item
// @Tags TripImages
// @Accept json
// @Produce json
// @Param tripId path int true "Trip ID"
// @Param tripItemId path int true "Trip Item ID"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.TripImageResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /trips/{tripId}/trip-items/{tripItemId}/images [get]
func (h *TripImageHandler) GetAllByTripIDAndTripItemID(c *gin.Context) {
	userID := middleware.GetUserIdHelper(c)

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		c.JSON(statusCode, errResponse)
		return
	}

	tripItemID, err := strconv.ParseInt(c.Param("tripItemId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripItemId")
		c.JSON(statusCode, errResponse)
		return
	}

	tripImages, errCode := h.tripImageService.GetAllByTripIDAndTripItemID(c, userID, tripID, tripItemID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		c.JSON(statusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&tripImages))
}
