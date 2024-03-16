// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/auth/check_user": {
            "get": {
                "description": "Get user info if user is authorized",
                "tags": [
                    "auth"
                ],
                "summary": "Check user",
                "operationId": "check-user",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/auth/login": {
            "post": {
                "description": "Login as a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in",
                "operationId": "sign-in",
                "parameters": [
                    {
                        "description": "login data",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserFormData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/logout": {
            "delete": {
                "description": "Quit from user` + "`" + `s account",
                "tags": [
                    "auth"
                ],
                "summary": "Log out",
                "operationId": "log-out",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/api/auth/signup": {
            "post": {
                "description": "Add a new user to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up",
                "operationId": "sign-up",
                "parameters": [
                    {
                        "description": "registration data",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserFormData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/note/add": {
            "post": {
                "description": "Create new note to current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "note"
                ],
                "summary": "Add note",
                "operationId": "add-note",
                "parameters": [
                    {
                        "description": "note data",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpsertNoteRequestForSwagger"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "note",
                        "schema": {
                            "$ref": "#/definitions/models.NoteForSwagger"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/note/get_all": {
            "get": {
                "description": "Get a list of notes of current user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "note"
                ],
                "summary": "Get all notes",
                "operationId": "get-all-notes",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "notes count",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "notes offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "notes title substring",
                        "name": "title",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "notes",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.NoteForSwagger"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/note/{id}": {
            "get": {
                "description": "Get one of notes of current user",
                "tags": [
                    "note"
                ],
                "summary": "Get one note",
                "operationId": "get-note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "note id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "note",
                        "schema": {
                            "$ref": "#/definitions/models.NoteForSwagger"
                        }
                    },
                    "400": {
                        "description": "incorrect id",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "note not found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/note/{id}/delete": {
            "delete": {
                "description": "Delete selected note of current user",
                "tags": [
                    "note"
                ],
                "summary": "Delete note",
                "operationId": "delete-note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "note id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "incorrect id",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "note not found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/note/{id}/edit": {
            "post": {
                "description": "Create new note to current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "note"
                ],
                "summary": "Update note",
                "operationId": "update-note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "note id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "note data",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpsertNoteRequestForSwagger"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "note",
                        "schema": {
                            "$ref": "#/definitions/models.NoteForSwagger"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/profile/get": {
            "get": {
                "description": "Get user info if user is authorized",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Get profile",
                "operationId": "get-profile",
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/profile/update": {
            "post": {
                "description": "Change password and/or description",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Update profile",
                "operationId": "update-profile",
                "parameters": [
                    {
                        "description": "update data",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ProfileUpdatePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/profile/update_avatar": {
            "post": {
                "description": "Change images",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Update profile images",
                "operationId": "update-profile-images",
                "parameters": [
                    {
                        "type": "file",
                        "description": "avatar",
                        "name": "avatar",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "413": {
                        "description": "Request Entity Too Large"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.NoteDataForSwagger": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.NoteForSwagger": {
            "type": "object",
            "properties": {
                "create_time": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/models.NoteDataForSwagger"
                },
                "id": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "update_time": {
                    "type": "string"
                }
            }
        },
        "models.ProfileUpdatePayload": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "password": {
                    "$ref": "#/definitions/models.passwords"
                }
            }
        },
        "models.UpsertNoteRequestForSwagger": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.NoteDataForSwagger"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "create_time": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image_path": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.UserFormData": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.passwords": {
            "type": "object",
            "properties": {
                "new": {
                    "type": "string"
                },
                "old": {
                    "type": "string"
                }
            }
        },
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "you-note.ru:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "YouNote API",
	Description:      "API for YouNote service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
