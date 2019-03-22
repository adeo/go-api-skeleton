package handlers

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/adeo/turbine-go-api-skeleton/middlewares"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao"
	dbFake "github.com/adeo/turbine-go-api-skeleton/storage/dao/fake"
	dbMock "github.com/adeo/turbine-go-api-skeleton/storage/dao/mock"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao/mongodb"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao/postgresql"
	"github.com/adeo/turbine-go-api-skeleton/storage/validators"
	"github.com/adeo/turbine-go-api-skeleton/utils"
	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ApplicationName      = ""
	ApplicationVersion   = "dev"
	ApplicationGitHash   = ""
	ApplicationBuildDate = ""
)

type Config struct {
	Mock                 bool
	DBInMemory           bool
	DBInMemoryImportFile string
	DBConnectionURI      string
	DBName               string
	PortAPI              int
	PortMonitoring       int
	LogLevel             string
	LogFormat            string
}

type Context struct {
	db        dao.Database
	validator *validator.Validate
}

func NewHandlersContext(config *Config) *Context {
	hc := &Context{}
	if config.Mock {
		hc.db = dbMock.NewDatabaseMock()
	} else if config.DBInMemory {
		hc.db = dbFake.NewDatabaseFake(config.DBInMemoryImportFile)
	} else if strings.HasPrefix(config.DBConnectionURI, "postgresql://") {
		hc.db = postgresql.NewDatabasePostgreSQL(config.DBConnectionURI)
	} else if strings.HasPrefix(config.DBConnectionURI, "mongodb://") {
		hc.db = mongodb.NewDatabaseMongoDB(config.DBConnectionURI, config.DBName)
	} else {
		utils.GetLogger().Fatal("no db connection uri given or not handled, starting in mode db in memory")
		hc.db = dbFake.NewDatabaseFake(config.DBInMemoryImportFile)
	}
	hc.validator = newValidator()
	return hc
}

func NewMonitoringRouter(hc *Context) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middlewares.GetLoggerMiddleware())
	router.Use(middlewares.GetHTTPLoggerMiddleware())

	public := router.Group("/")
	public.Use(middlewares.CORSMiddlewareForOthersHTTPMethods())

	public.Handle(http.MethodGet, "/info", hc.GetInfo)
	public.Handle(http.MethodOptions, "/info", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))
	public.Handle(http.MethodGet, "/openapi", hc.GetOpenAPISchema)
	public.Handle(http.MethodOptions, "/openapi", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))
	public.Handle(http.MethodGet, "/prometheus", gin.WrapH(promhttp.Handler()))
	public.Handle(http.MethodOptions, "/prometheus", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))

	if dbInMemory, ok := hc.db.(*dbFake.DatabaseFake); ok {
		// db in memory mode, add export endpoint
		public.Handle(http.MethodGet, "/export", func(c *gin.Context) {
			httputils.JSON(c.Writer, http.StatusOK, dbInMemory.Export())
		})
	}

	return router
}

func NewAPIRouter(hc *Context) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middlewares.NewPrometheusMiddleware(ApplicationName).Handler())
	router.Use(middlewares.GetLoggerMiddleware())
	router.Use(middlewares.GetHTTPLoggerMiddleware())

	public := router.Group("/")
	public.Use(middlewares.CORSMiddlewareForOthersHTTPMethods())

	// start: user routes
	public.Handle(http.MethodOptions, "/users", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPost))
	public.Handle(http.MethodGet, "/users", hc.GetAllUsers)
	public.Handle(http.MethodPost, "/users", hc.CreateUser)
	public.Handle(http.MethodOptions, "/users/:id", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPut, http.MethodDelete))
	public.Handle(http.MethodGet, "/users/:id", hc.GetUser)
	public.Handle(http.MethodPut, "/users/:id", hc.UpdateUser)
	public.Handle(http.MethodDelete, "/users/:id", hc.DeleteUser)
	// end: user routes
	// start: template routes
	public.Handle(http.MethodOptions, "/templates", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPost))
	public.Handle(http.MethodGet, "/templates", hc.GetAllTemplates)
	public.Handle(http.MethodPost, "/templates", hc.CreateTemplate)
	public.Handle(http.MethodOptions, "/templates/:id", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPut, http.MethodDelete))
	public.Handle(http.MethodGet, "/templates/:id", hc.GetTemplate)
	public.Handle(http.MethodPut, "/templates/:id", hc.UpdateTemplate)
	public.Handle(http.MethodDelete, "/templates/:id", hc.DeleteTemplate)
	// end: template routes

	return router
}

func newValidator() *validator.Validate {
	va := validator.New()

	va.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(name) < 1 {
			return ""
		}
		return name[0]
	})

	for k, v := range validators.CustomValidators {
		if v.Validator != nil {
			va.RegisterValidationCtx(k, v.Validator)
		}
	}

	return va
}

func (hc *Context) getValidationContext(c *gin.Context) context.Context {
	vc := &validators.ValidationContext{
		DB: hc.db,
	}
	return context.WithValue(c, validators.ContextKeyValidator, vc)
}
