package handlers

import (
	"net/http"
	"strings"

	authentication "github.com/adeo/turbine-auth/pkg/client/http"
	"github.com/adeo/turbine-auth/pkg/client/middleware"
	"github.com/adeo/turbine-go-api-skeleton/middlewares"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao"
	dbFake "github.com/adeo/turbine-go-api-skeleton/storage/dao/fake" // DAO IN MEMORY
	dbMock "github.com/adeo/turbine-go-api-skeleton/storage/dao/mock"
	"github.com/adeo/turbine-go-api-skeleton/storage/dao/mongodb"    // DAO MONGO
	"github.com/adeo/turbine-go-api-skeleton/storage/dao/postgresql" // DAO PG
	"github.com/adeo/turbine-go-api-skeleton/storage/validators"
	"github.com/adeo/turbine-go-api-skeleton/utils"
	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/go-playground/validator.v9"
)

const (
	baseURI = ""
)

var (
	ApplicationName      = ""
	ApplicationVersion   = "dev"
	ApplicationGitHash   = ""
	ApplicationBuildDate = ""
)

type Config struct {
	Mock                      bool
	DBInMemory                bool   // DAO IN MEMORY
	DBInMemoryImportFile      string // DAO IN MEMORY
	DBConnectionURI           string
	DBName                    string
	PortAPI                   int
	PortMonitoring            int
	LogLevel                  string
	LogFormat                 string
	AuthenticationServiceFake bool
	AuthenticationServiceURI  string
}

type Context struct {
	db                    dao.Database
	authenticationService authentication.Service
	validator             *validator.Validate
}

func NewHandlersContext(config *Config) *Context {
	hc := &Context{}

	if config.Mock {
		hc.db = dbMock.NewDatabaseMock()
	} else if config.DBInMemory { // DAO IN MEMORY
		hc.db = dbFake.NewDatabaseFake(config.DBInMemoryImportFile) // DAO IN MEMORY
	} else if strings.HasPrefix(config.DBConnectionURI, "postgresql://") { // DAO PG
		hc.db = postgresql.NewDatabasePostgreSQL(config.DBConnectionURI) // DAO PG
	} else if strings.HasPrefix(config.DBConnectionURI, "mongodb://") { // DAO MONGO
		hc.db = mongodb.NewDatabaseMongoDB(config.DBConnectionURI, config.DBName) // DAO MONGO
	} else {
		utils.GetLogger().Fatal("no db connection uri given or not handled, and no db in memory mode enabled, exiting")
	}

	if config.AuthenticationServiceFake {
		hc.authenticationService = authentication.NewServiceFake()
	} else {
		hc.authenticationService = authentication.NewServiceHTTP(config.AuthenticationServiceURI)
	}

	hc.validator = validators.NewValidator()
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
	public.Use(middlewares.GetCORSMiddlewareForOthersHTTPMethods())

	public.Handle(http.MethodGet, "/_health", hc.GetHealth)
	public.Handle(http.MethodOptions, "/_health", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))

	public.Handle(http.MethodGet, "/openapi", hc.GetOpenAPISchema)
	public.Handle(http.MethodOptions, "/openapi", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))

	public.Handle(http.MethodGet, "/prometheus", gin.WrapH(promhttp.Handler()))
	public.Handle(http.MethodOptions, "/prometheus", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet))

	if dbInMemory, ok := hc.db.(*dbFake.DatabaseFake); ok { // DAO IN MEMORY
		// db in memory mode, add export endpoint // DAO IN MEMORY
		public.Handle(http.MethodGet, "/export", func(c *gin.Context) { // DAO IN MEMORY
			httputils.JSON(c.Writer, http.StatusOK, dbInMemory.Export()) // DAO IN MEMORY
		}) // DAO IN MEMORY
	} // DAO IN MEMORY

	return router
}

func NewAPIRouter(hc *Context) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(gin.Recovery())
	router.Use(middlewares.GetPrometheusMiddleware(ApplicationName))
	router.Use(middlewares.GetLoggerMiddleware())
	router.Use(middlewares.GetHTTPLoggerMiddleware())

	handleAPIRoutes(hc, router)
	handleCORSRoutes(hc, router)

	return router
}

func handleCORSRoutes(hc *Context, router *gin.Engine) {
	public := router.Group(baseURI)

	// start: template routes
	public.Handle(http.MethodOptions, "/templates", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPost))
	public.Handle(http.MethodOptions, "/templates/:id", hc.GetOptionsHandler(httputils.AllowedHeaders, http.MethodGet, http.MethodPut, http.MethodDelete))
	// end: template routes
}

func handleAPIRoutes(hc *Context, router *gin.Engine) {
	public := router.Group(baseURI)
	public.Use(middlewares.GetCORSMiddlewareForOthersHTTPMethods())

	secured := public.Group("/")
	secured.Use(middleware.GetAuthenticationMiddleware(hc.authenticationService))

	// start: template routes
	secured.Handle(http.MethodGet, "/templates", hc.GetAllTemplates)
	secured.Handle(http.MethodPost, "/templates", hc.CreateTemplate)
	secured.Handle(http.MethodGet, "/templates/:id", hc.GetTemplate)
	secured.Handle(http.MethodPut, "/templates/:id", hc.UpdateTemplate)
	secured.Handle(http.MethodDelete, "/templates/:id", hc.DeleteTemplate)
	// end: template routes
}
