// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tasks": {
            "get": {
                "description": "Get a list of tasks with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "List tasks",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "Filter by completion status",
                        "name": "completed",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "date",
                        "example": "2024-06-07",
                        "description": "Filter by due date",
                        "name": "date",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of tasks per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of tasks",
                        "schema": {
                            "$ref": "#/definitions/todo.Pages"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new task with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Create a new task",
                "parameters": [
                    {
                        "description": "Task to create",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/todo.Task"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Task created successfully",
                        "schema": {
                            "$ref": "#/definitions/todo.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tasks/{id}": {
            "get": {
                "description": "Get a task by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get a task by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task found",
                        "schema": {
                            "$ref": "#/definitions/todo.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a task with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Update a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task updates",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task updated successfully",
                        "schema": {
                            "$ref": "#/definitions/todo.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a task by ID",
                "tags": [
                    "tasks"
                ],
                "summary": "Delete a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/todo.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "todo.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Error"
                }
            }
        },
        "todo.Pages": {
            "type": "object",
            "properties": {
                "count_page": {
                    "type": "integer"
                },
                "cur_page": {
                    "type": "integer"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/todo.Task"
                    }
                }
            }
        },
        "todo.Task": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "due_date": {
                    "type": "string",
                    "example": "2024-06-07T15:00:00Z"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
