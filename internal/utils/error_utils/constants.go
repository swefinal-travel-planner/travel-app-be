package error_utils

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type systemErrorMessage struct {
	SqlxNoRow       string
	RedisNil        string
	NotMemberOfTrip string
}

var SystemErrorMessage = systemErrorMessage{
	SqlxNoRow:       sql.ErrNoRows.Error(),
	RedisNil:        redis.Nil.Error(),
	NotMemberOfTrip: "not a member of the trip",
}

type errorCode struct {
	// redis related
	REDIS_DOWN string

	// db related
	DB_DOWN string

	// auth related
	REGISTER_INVALID_OTP                       string
	REGISTER_EMAIL_EXISTED                     string
	LOGIN_EMAIL_NOT_FOUND                      string
	LOGIN_INVALID_PASSWORD                     string
	REGISTER_SEND_OTP_TO_EXISTED_EMAIL         string
	REGISTER_OTP_NOT_FOUND                     string
	REGISTER_OTP_INVALID                       string
	RESET_PASSWORD_EMAIL_NOT_FOUND             string
	RESET_PASSWORD_OTP_NOT_FOUND               string
	RESET_PASSWORD_OTP_INVALID                 string
	SET_PASSWORD_OTP_INVALID                   string
	SET_PASSWORD_OTP_NOT_FOUND                 string
	REFRESH_TOKEN_NOT_FOUND                    string
	REFRESH_TOKEN_INVALID                      string
	ACCESS_TOKEN_INVALID                       string
	REMOVE_FRIEND_NOT_FOUND                    string
	ADD_FRIEND_RECEIVER_NOT_FOUND              string
	ADD_FRIEND_CANNOT_ADD_YOURSELF             string
	ADD_FRIEND_IN_COOLDOWN                     string
	ADD_FRIEND_INVITATION_ALREADY_EXISTS       string
	ADD_FRIEND_ALREADY_FRIEND                  string
	FRIEND_INVITATION_NOT_FOUND                string
	FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER  string
	FRIEND_INVITATION_ONLY_SENDER_CAN_WITHDRAW string
	TRIP_NOT_FOUND                             string
	FORBIDDEN                                  string

	TRIP_INVITATION_ALREADY_EXISTS          string
	TRIP_INVITATION_RECEIVER_ALREADY_MEMBER string
	INTERNAL_SERVER_ERROR                   string
	BAD_REQUEST                             string
}

var ErrorCode = errorCode{
	REDIS_DOWN: "REDIS_DOWN",
	DB_DOWN:    "DB_DOWN",

	REGISTER_EMAIL_EXISTED:                     "REGISTER_EMAIL_EXISTED",
	LOGIN_EMAIL_NOT_FOUND:                      "LOGIN_EMAIL_NOT_FOUND",
	LOGIN_INVALID_PASSWORD:                     "LOGIN_INVALID_PASSWORD",
	REGISTER_SEND_OTP_TO_EXISTED_EMAIL:         "REGISTER_SEND_OTP_TO_EXISTED_EMAIL",
	REGISTER_OTP_NOT_FOUND:                     "REGISTER_OTP_NOT_FOUND",
	REGISTER_OTP_INVALID:                       "REGISTER_OTP_INVALID",
	RESET_PASSWORD_EMAIL_NOT_FOUND:             "RESET_PASSWORD_EMAIL_NOT_FOUND",
	RESET_PASSWORD_OTP_NOT_FOUND:               "RESET_PASSWORD_OTP_NOT_FOUND",
	RESET_PASSWORD_OTP_INVALID:                 "RESET_PASSWORD_OTP_INVALID",
	SET_PASSWORD_OTP_INVALID:                   "SET_PASSWORD_OTP_INVALID",
	SET_PASSWORD_OTP_NOT_FOUND:                 "SET_PASSWORD_OTP_NOT_FOUND",
	REFRESH_TOKEN_INVALID:                      "REFRESH_TOKEN_INVALID",
	ACCESS_TOKEN_INVALID:                       "ACCESS_TOKEN_INVALID",
	REMOVE_FRIEND_NOT_FOUND:                    "REMOVE_FRIEND_NOT_FOUND",
	ADD_FRIEND_RECEIVER_NOT_FOUND:              "ADD_FRIEND_RECEIVER_NOT_FOUND",
	ADD_FRIEND_CANNOT_ADD_YOURSELF:             "ADD_FRIEND_CANNOT_ADD_YOURSELF",
	ADD_FRIEND_IN_COOLDOWN:                     "ADD_FRIEND_IN_COOLDOWN",
	ADD_FRIEND_INVITATION_ALREADY_EXISTS:       "ADD_FRIEND_INVITATION_ALREADY_EXISTS",
	ADD_FRIEND_ALREADY_FRIEND:                  "ADD_FRIEND_ALREADY_FRIEND",
	FRIEND_INVITATION_NOT_FOUND:                "FRIEND_INVITATION_NOT_FOUND",
	FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER:  "FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER",
	FRIEND_INVITATION_ONLY_SENDER_CAN_WITHDRAW: "FRIEND_INVITATION_CAN_WITHDRAW",
	TRIP_NOT_FOUND:                             "TRIP_NOT_FOUND",
	FORBIDDEN:                                  "FORBIDDEN",

	TRIP_INVITATION_ALREADY_EXISTS:          "TRIP_INVITATION_ALREADY_EXISTS",
	TRIP_INVITATION_RECEIVER_ALREADY_MEMBER: "TRIP_INVITATION_RECEIVER_ALREADY_MEMBER",
	INTERNAL_SERVER_ERROR:                   "INTERNAL_SERVER_ERROR",
	BAD_REQUEST:                             "BAD_REQUEST",
}
