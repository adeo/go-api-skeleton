package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adeo/turbine-go-api-skeleton/middlewares"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao"
	"github.com/adeo/turbine-go-api-skeleton/storage/model"
	"github.com/adeo/turbine-go-api-skeleton/storage/validators"
	"github.com/adeo/turbine-go-api-skeleton/utils"
	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

// @openapi:path
// /templates:
//	get:
//		tags:
//			- templates
//		description: "Get all the templates"
//		responses:
//			200:
//				description: "The array containing the templates"
//				content:
//					application/json:
//						schema:
//							type: "array"
//							items:
//								$ref: "#/components/schemas/Template"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) GetAllTemplates(c *gin.Context) {
	templates, err := hc.db.GetAllTemplates()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while getting templates")
		httputils.JSONErrorWithMessage(c.Writer, model.ErrInternalServer, "Error while getting templates")
		return
	}
	httputils.JSONOK(c, templates)
}

// @openapi:path
// /templates:
//	post:
//		tags:
//			- templates
//		description: "Create a new template"
//		requestBody:
//			description: The template data.
//			required: true
//			content:
//				application/json:
//					schema:
//						$ref: "#/components/schemas/TemplateEditable"
//		responses:
//			201:
//				description: "The created template"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			400:
//				description: "This error occurs when the request is not correct (bad body format, validation error)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			409:
//				description: "This error occurs when the new entity is in conflict with exiting one (duplicated)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) CreateTemplate(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating template, read data fail")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	templateToCreate := model.TemplateEditable{}
	err = json.Unmarshal(body, &templateToCreate)
	if err != nil {
		httputils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(validators.NewContextWithValidationContext(c, hc.db), templateToCreate)
	if err != nil {
		httputils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template := &model.Template{
		TemplateEditable: templateToCreate,
	}

	err = hc.db.CreateTemplate(template)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeDuplicate:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrAlreadyExists, "Template already exists")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error CreateTemplate: Error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while creating template")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	c.Writer.Header().Set(httputils.HeaderNameLocation, fmt.Sprintf("%s/templates/%s", baseURI, template.Name))
	httputils.JSON(c.Writer, http.StatusCreated, template)
}

// @openapi:path
// /templates/{templateID}:
//	get:
//		tags:
//			- templates
//		description: "Get a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to get"
//		responses:
//			200:
//				description: "The templates with id `templateID`"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) GetTemplate(c *gin.Context) {
	c.Set(middlewares.ContextKeyPrometheusURI, baseURI+"/templates/:id")

	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
	if err != nil {
		httputils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template, err := hc.db.GetTemplateByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error GetTemplate: get template error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	if template == nil {
		httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template not found")
		return
	}

	httputils.JSONOK(c, template)
}

// @openapi:path
// /templates/{templateID}:
//	delete:
//		tags:
//			- templates
//		description: "Delete a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to delete"
//		responses:
//			204:
//				description: "Templates with id `templateID` deleted"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) DeleteTemplate(c *gin.Context) {
	c.Set(middlewares.ContextKeyPrometheusURI, baseURI+"/templates/:id")

	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
	if err != nil {
		httputils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check template id given in URL exists
	_, err = hc.db.GetTemplateByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to delete not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: get template error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template to delete")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	err = hc.db.DeleteTemplate(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to delete not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error DeleteTemplate: Error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while deleting template")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	httputils.JSON(c.Writer, http.StatusNoContent, nil)
}

// @openapi:path
// /templates/{templateID}:
//	put:
//		tags:
//			- templates
//		description: "Update a template"
//		parameters:
//		- in: path
//		  name: templateID
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template id to update"
//		- in: header
//		  name: If-Match
//		  schema:
//		  	type: string
//		  required: true
//		  description: "The template version to update. You can find the template version using the GET endpoint, in the ETag response header. If the version has been updated between your GET and your PUT, you will receive a 412 Precondition Failed response."
//		requestBody:
//			description: The template data.
//			required: true
//			content:
//				application/json:
//					schema:
//						$ref: "#/components/schemas/TemplateEditable"
//		responses:
//			200:
//				description: "The updated template"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/Template"
//			400:
//				description: "This error occurs when the request is not correct (bad body format, validation error)"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			404:
//				description: "Template not found"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
//			500:
//				description: "Server error"
//				content:
//					application/json:
//						schema:
//							$ref: "#/components/schemas/APIError"
func (hc *Context) UpdateTemplate(c *gin.Context) {
	c.Set(middlewares.ContextKeyPrometheusURI, baseURI+"/templates/:id")

	templateID := c.Param("id")

	err := hc.validator.VarCtx(c, templateID, "required")
	if err != nil {
		httputils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	// check template id given in URL exists
	template, err := hc.db.GetTemplateByID(templateID)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to update not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("UpdateTemplate: get template error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while get template to update")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	// check versions
	if !utils.IsSameVersion(c.GetHeader(httputils.HeaderNameIfMatch), template) {
		httputils.JSONError(c.Writer, model.ErrVersionMismatched)
		return
	}

	// get body and verify data
	body, err := c.GetRawData()
	if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating template, read data fail")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	templateToUpdate := model.TemplateEditable{}
	err = json.Unmarshal(body, &templateToUpdate)
	if err != nil {
		httputils.JSONError(c.Writer, model.ErrBadRequestFormat)
		return
	}

	err = hc.validator.StructCtx(validators.NewContextWithValidationContext(c, hc.db), templateToUpdate)
	if err != nil {
		httputils.JSONError(c.Writer, validators.NewDataValidationAPIError(err))
		return
	}

	template.TemplateEditable = templateToUpdate

	// make the update
	err = hc.db.UpdateTemplate(template)
	if e, ok := err.(*dao.DAOError); ok {
		switch {
		case e.Type == dao.ErrTypeNotFound:
			httputils.JSONErrorWithMessage(c.Writer, model.ErrNotFound, "Template to update not found")
			return
		default:
			utils.GetLoggerFromCtx(c).WithError(err).WithField("type", e.Type).Error("error UpdateTemplate: Error type not handled")
			httputils.JSONError(c.Writer, model.ErrInternalServer)
			return
		}
	} else if err != nil {
		utils.GetLoggerFromCtx(c).WithError(err).Error("error while updating template")
		httputils.JSONError(c.Writer, model.ErrInternalServer)
		return
	}

	httputils.JSON(c.Writer, http.StatusOK, template)
}
