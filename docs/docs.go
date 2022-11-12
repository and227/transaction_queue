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
        "/transactions": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transaction"
                ],
                "summary": "create new transaction",
                "operationId": "create-transaction",
                "parameters": [
                    {
                        "description": "transaction info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.TransactionCreateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Transaction"
                        }
                    },
                    "400": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "500": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "create new user",
                "operationId": "create-user",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.UserCreateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.User"
                        }
                    },
                    "400": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "500": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        },
        "/users/{user_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get user with balance",
                "operationId": "get-user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.UserWithBalanceOutDTO"
                        }
                    },
                    "400": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "500": {
                        "description": "error response",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.BalanceOutDTO": {
            "type": "object",
            "required": [
                "amount",
                "id"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "db.Transaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "db.TransactionCreateDTO": {
            "type": "object",
            "required": [
                "amount",
                "type",
                "user_id"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "db.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "db.UserCreateDTO": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "db.UserWithBalanceOutDTO": {
            "type": "object",
            "required": [
                "id",
                "name"
            ],
            "properties": {
                "balance": {
                    "$ref": "#/definitions/db.BalanceOutDTO"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "gin.H": {
            "type": "object",
            "additionalProperties": {}
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Simple transaction queue",
	Description:      "Simple transaction queue",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}