package docs

import "github.com/swaggo/swag"

// docTemplate is a lightweight OpenAPI spec for Swagger UI.
const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "Restaurant Management System API",
        "title": "RMS API",
        "version": "1.0.0"
    },
    "basePath": "/",
    "schemes": ["http"],
    "paths": {}
}`

type s struct{}

// ReadDoc returns the template spec to Swagger middleware.
func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register("swagger", &s{})
}
