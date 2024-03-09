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
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Check user",
                "operationId": "check-user",
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
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
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
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/note/add": {
            "post": {
                "description": "Create new note to current user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "note"
                ],
                "summary": "Add note",
                "operationId": "add-note",
                "responses": {
                    "200": {
                        "description": "note",
                        "schema": {
                            "$ref": "#/definitions/models.Note"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
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
                                "$ref": "#/definitions/models.Note"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
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
                "produces": [
                    "application/json"
                ],
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
                            "$ref": "#/definitions/models.Note"
                        }
                    },
                    "400": {
                        "description": "incorrect id",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "note not found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Note": {
            "type": "object",
            "properties": {
                "create_time": {
                    "type": "string"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
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
        "utils.ErrorResponse": {
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
