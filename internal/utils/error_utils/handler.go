package error_utils

import (
	"net/http"

	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
)

func ErrorCodeToHttpResponse(errCode string, field string) (statusCode int, httpErrResponse httpcommon.HttpResponse[any]) {
	switch errCode {
	case ErrorCode.INTERNAL_SERVER_ERROR:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "An unexpected error occurred. Please try again later or contact support if the problem persists",
			Field:   field,
			Code:    ErrorCode.INTERNAL_SERVER_ERROR,
		})
	case ErrorCode.DB_DOWN:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Our database system is currently unavailable. Please try again in a few minutes",
			Field:   field,
			Code:    ErrorCode.DB_DOWN,
		})
	case ErrorCode.REDIS_DOWN:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Our authentication system is temporarily unavailable. OTP verification and related functions will not work. Please try again later",
			Field:   field,
			Code:    ErrorCode.REDIS_DOWN,
		})
	case ErrorCode.REGISTER_INVALID_OTP:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The OTP code you entered is invalid. Please check your email and enter the correct code",
			Field:   field,
			Code:    ErrorCode.REGISTER_INVALID_OTP,
		})
	case ErrorCode.REGISTER_EMAIL_EXISTED:
		statusCode = http.StatusConflict
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "This email address is already registered. Please use a different email or try logging in",
			Field:   field,
			Code:    ErrorCode.REGISTER_EMAIL_EXISTED,
		})
	case ErrorCode.LOGIN_EMAIL_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "No account found with this email address. Please check your email or register a new account",
			Field:   field,
			Code:    ErrorCode.LOGIN_EMAIL_NOT_FOUND,
		})
	case ErrorCode.LOGIN_INVALID_PASSWORD:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The password you entered is incorrect. Please try again or use the 'Forgot Password' option",
			Field:   field,
			Code:    ErrorCode.LOGIN_INVALID_PASSWORD,
		})
	case ErrorCode.REGISTER_SEND_OTP_TO_EXISTED_EMAIL:
		statusCode = http.StatusConflict
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "This email is already registered. OTP cannot be sent to existing accounts. Please try logging in instead",
			Field:   field,
			Code:    ErrorCode.REGISTER_SEND_OTP_TO_EXISTED_EMAIL,
		})
	case ErrorCode.REGISTER_OTP_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The OTP code has expired or was not found. Please request a new OTP code",
			Field:   field,
			Code:    ErrorCode.REGISTER_OTP_NOT_FOUND,
		})
	case ErrorCode.REGISTER_OTP_INVALID:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The OTP code you entered is incorrect. Please check your email and try again",
			Field:   field,
			Code:    ErrorCode.REGISTER_OTP_INVALID,
		})
	case ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "No account found with this email address. Please check your email or register a new account",
			Field:   field,
			Code:    ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND,
		})
	case ErrorCode.RESET_PASSWORD_OTP_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The password reset OTP has expired or was not found. Please request a new password reset code",
			Field:   field,
			Code:    ErrorCode.RESET_PASSWORD_OTP_NOT_FOUND,
		})
	case ErrorCode.RESET_PASSWORD_OTP_INVALID:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The password reset code you entered is incorrect. Please check your email and try again",
			Field:   field,
			Code:    ErrorCode.RESET_PASSWORD_OTP_INVALID,
		})
	case ErrorCode.SET_PASSWORD_OTP_INVALID:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The password setup code you entered is incorrect. Please check your email and try again",
			Field:   field,
			Code:    ErrorCode.SET_PASSWORD_OTP_INVALID,
		})
	case ErrorCode.SET_PASSWORD_OTP_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The password setup code has expired or was not found. Please request a new setup code",
			Field:   field,
			Code:    ErrorCode.SET_PASSWORD_OTP_NOT_FOUND,
		})
	case ErrorCode.REFRESH_TOKEN_NOT_FOUND:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Your session has expired. Please log in again to continue",
			Field:   field,
			Code:    ErrorCode.REFRESH_TOKEN_NOT_FOUND,
		})
	case ErrorCode.REFRESH_TOKEN_INVALID:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Your session is invalid. Please log in again to continue",
			Field:   field,
			Code:    ErrorCode.REFRESH_TOKEN_INVALID,
		})
	case ErrorCode.ACCESS_TOKEN_INVALID:
		statusCode = http.StatusUnauthorized
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Your access token is no longer valid",
			Field:   field,
			Code:    ErrorCode.ACCESS_TOKEN_INVALID,
		})
	case ErrorCode.REMOVE_FRIEND_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Cannot remove friend: The friendship between these users does not exist",
			Field:   field,
			Code:    ErrorCode.REMOVE_FRIEND_NOT_FOUND,
		})
	case ErrorCode.BAD_REQUEST:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Invalid request parameters",
			Field:   field,
			Code:    ErrorCode.BAD_REQUEST,
		})
	case ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The user you are trying to add as friend does not exist",
			Field:   field,
			Code:    ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND,
		})
	case ErrorCode.ADD_FRIEND_CANNOT_ADD_YOURSELF:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "You cannot add yourself as a friend",
			Field:   field,
			Code:    ErrorCode.ADD_FRIEND_CANNOT_ADD_YOURSELF,
		})
	case ErrorCode.ADD_FRIEND_IN_COOLDOWN:
		statusCode = http.StatusTooManyRequests
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Friend request to this user has been cooled down. Please wait before sending another one",
			Field:   field,
			Code:    ErrorCode.ADD_FRIEND_IN_COOLDOWN,
		})
	case ErrorCode.ADD_FRIEND_INVITATION_ALREADY_EXISTS:
		statusCode = http.StatusConflict
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "A friend invitation already exists between you and this user",
			Field:   field,
			Code:    ErrorCode.ADD_FRIEND_INVITATION_ALREADY_EXISTS,
		})
	case ErrorCode.ADD_FRIEND_ALREADY_FRIEND:
		statusCode = http.StatusConflict
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "You are already friends with this user",
			Field:   field,
			Code:    ErrorCode.ADD_FRIEND_ALREADY_FRIEND,
		})
	case ErrorCode.FRIEND_INVITATION_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The friend invitation does not exist",
			Field:   field,
			Code:    ErrorCode.FRIEND_INVITATION_NOT_FOUND,
		})
	case ErrorCode.FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER:
		statusCode = http.StatusBadRequest
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "You cannot accept a friend invitation that you sent",
			Field:   field,
			Code:    ErrorCode.FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER,
		})
	case ErrorCode.FRIEND_INVITATION_ONLY_SENDER_CAN_WITHDRAW:
		statusCode = http.StatusForbidden
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Only the sender can withdraw a friend invitation",
			Field:   field,
			Code:    ErrorCode.FRIEND_INVITATION_ONLY_SENDER_CAN_WITHDRAW,
		})
	case ErrorCode.TRIP_NOT_FOUND:
		statusCode = http.StatusNotFound
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "The trip you are looking for does not exist",
			Field:   field,
			Code:    ErrorCode.TRIP_NOT_FOUND,
		})
	case ErrorCode.FORBIDDEN:
		statusCode = http.StatusForbidden
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "You do not have permission to perform this action",
			Field:   field,
			Code:    ErrorCode.FORBIDDEN,
		})
	default:
		statusCode = http.StatusInternalServerError
		httpErrResponse = httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "An unexpected error occurred. Please try again later or contact support if the problem persists",
			Field:   field,
			Code:    ErrorCode.INTERNAL_SERVER_ERROR,
		})
	}

	return
}
