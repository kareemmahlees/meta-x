// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Kareem Ebrahim",
            "email": "kareemmahlees@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "get info about the api",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.APIInfoResult"
                        }
                    }
                }
            }
        },
        "/databases": {
            "get": {
                "description": "list databases",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Databases"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListDatabasesResult"
                        }
                    }
                }
            },
            "post": {
                "description": "create pg/mysql database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Databases"
                ],
                "parameters": [
                    {
                        "description": "only supported for pg and mysql, because attached sqlite dbs are temporary",
                        "name": "pg_mysql_db_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreatePgMySqlDBPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CreateDatabaseResult"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "check application health by getting current date",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HealthCheckResult"
                        }
                    }
                }
            }
        },
        "/tables": {
            "get": {
                "description": "list tables",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tables"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HandleListTablesResp"
                        }
                    }
                }
            }
        },
        "/tables/{tableName}": {
            "put": {
                "description": "update table",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tables"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "table name",
                        "name": "tableName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update table data",
                        "name": "tableData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/lib.UpdateTableProps"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HandleUpdateDeleteResp"
                        }
                    }
                }
            },
            "post": {
                "description": "create table",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tables"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "table name",
                        "name": "tableName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "create table data",
                        "name": "tableData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.HandleCreateTableBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/routes.HandleCreateTableResp"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tables"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "table name",
                        "name": "tableName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HandleUpdateDeleteResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "lib.CreateTableProps": {
            "type": "object",
            "required": [
                "type"
            ],
            "properties": {
                "default": {},
                "nullable": {},
                "type": {
                    "type": "string",
                    "enum": [
                        "text",
                        "number",
                        "date"
                    ]
                },
                "unique": {}
            }
        },
        "lib.UpdateTableProps": {
            "type": "object",
            "required": [
                "operation"
            ],
            "properties": {
                "operation": {
                    "type": "object",
                    "required": [
                        "data",
                        "type"
                    ],
                    "properties": {
                        "data": {},
                        "type": {
                            "type": "string",
                            "enum": [
                                "add",
                                "modify",
                                "delete"
                            ]
                        }
                    }
                }
            }
        },
        "models.CreateDatabaseResult": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "integer"
                }
            }
        },
        "models.CreatePgMySqlDBPayload": {
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
        "models.ListDatabasesResult": {
            "type": "object",
            "properties": {
                "databases": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "routes.APIInfoResult": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "contact": {
                    "type": "string"
                },
                "repo": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "routes.HandleCreateTableBody": {
            "type": "object",
            "properties": {
                "colName": {
                    "$ref": "#/definitions/lib.CreateTableProps"
                }
            }
        },
        "routes.HandleCreateTableResp": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                }
            }
        },
        "routes.HandleListTablesResp": {
            "type": "object",
            "properties": {
                "tables": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "routes.HandleUpdateDeleteResp": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean"
                }
            }
        },
        "routes.HealthCheckResult": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5522",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "MySQL Meta",
	Description:      "A RESTFull and GraphQL API to manage your MySQL DB",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
