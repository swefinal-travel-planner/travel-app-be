// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/forgot-password": {
            "post": {
                "description": "Set a new password after OTP verification",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Set Password",
                "parameters": [
                    {
                        "description": "Set Password payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/forgot-password/otp": {
            "post": {
                "description": "Send OTP to user email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Send OTP to Mail",
                "parameters": [
                    {
                        "description": "Send OTP payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SendOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/forgot-password/verify-otp": {
            "post": {
                "description": "Verify OTP with email and otp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Verify OTP",
                "parameters": [
                    {
                        "description": "Verify OTP payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.VerifyOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/google-login": {
            "post": {
                "description": "Authenticate a user using their Google token via Firebase",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Login with Google",
                "parameters": [
                    {
                        "description": "Google Login Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.GoogleLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-entity_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Login to account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Auth payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-entity_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "Refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auths"
                ],
                "summary": "Refresh",
                "parameters": [
                    {
                        "description": "Auth payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-entity_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "User register",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Auth payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RegisterRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Authorization: Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpcommon.HttpResponse-any"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.User": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "googleId": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "httpcommon.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httpcommon.HttpResponse-any": {
            "type": "object",
            "properties": {
                "data": {},
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/httpcommon.Error"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "httpcommon.HttpResponse-entity_User": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/entity.User"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/httpcommon.Error"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "model.GoogleLoginRequest": {
            "type": "object",
            "required": [
                "displayName",
                "email",
                "password",
                "phoneNumber",
                "photoURL",
                "uid"
            ],
            "properties": {
                "displayName": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 5
                },
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "password": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 8
                },
                "phoneNumber": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "photoURL": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "uid": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                }
            }
        },
        "model.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "password": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 8
                }
            }
        },
        "model.RefreshTokenRequest": {
            "type": "object",
            "required": [
                "refreshToken"
            ],
            "properties": {
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "model.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "phoneNumber"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "name": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 5
                },
                "password": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 8
                },
                "phoneNumber": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                }
            }
        },
        "model.SendOTPRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                }
            }
        },
        "model.SetPasswordRequest": {
            "type": "object",
            "required": [
                "email",
                "otp",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "otp": {
                    "type": "string",
                    "maxLength": 6,
                    "minLength": 6
                },
                "password": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 8
                }
            }
        },
        "model.VerifyOTPRequest": {
            "type": "object",
            "required": [
                "email",
                "otp"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 10
                },
                "otp": {
                    "type": "string",
                    "maxLength": 6,
                    "minLength": 6
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
