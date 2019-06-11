package model

import "time"

// @openapi:schema
type Template struct {
	TemplateEditable `bson:",inline"` // avoid having a property "TemplateEditable" in your mongodb document
	ID               string           `json:"id" bson:"_id"`
	CreatedAt        time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt        *time.Time       `json:"updated_at" bson:"updated_at"`
}

// @openapi:schema
type TemplateEditable struct {
	// Add here your model properties, and don't forget to modify SQL request in corresponding DAO file if any
	Name string `json:"name" bson:"name" validate:"required"`
}
