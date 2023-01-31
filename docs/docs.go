// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/blogs/": {
            "get": {
                "description": "Show all blogs registered",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Blogs"
                ],
                "summary": "List all the Blogs",
                "operationId": "Fetch Blog",
                "parameters": [
                    {
                        "type": "string",
                        "description": "1",
                        "name": "Page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "10",
                        "name": "ItemsPerPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.listBlogsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/blogs/create-blog/": {
            "post": {
                "description": "Create a blog",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Blogs"
                ],
                "summary": "Create a blog",
                "operationId": "Create a blog",
                "parameters": [
                    {
                        "description": "Create blog request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.createBlogRequestBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/v1.createBlogResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/blogs/{id}": {
            "get": {
                "description": "Show a single blog registered",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Blogs"
                ],
                "summary": "Fetch single blog by ID",
                "operationId": "Single blog",
                "parameters": [
                    {
                        "type": "string",
                        "description": "blog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.singleBlogResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Blog": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "descriptions": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userRole": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "v1.createBlogRequestBody": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                }
            }
        },
        "v1.createBlogResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.listBlogsResponse": {
            "type": "object",
            "properties": {
                "blogs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.createBlogResponse"
                    }
                }
            }
        },
        "v1.singleBlogResponse": {
            "type": "object"
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8089",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Gin Gonic golang Clean Architecture.",
	Description: "Illustration of uncle Bob's clean architecture using a demo blogs api.\nIt serves as Blog.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}