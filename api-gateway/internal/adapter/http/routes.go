package http

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	_ "github.com/sorawaslocked/car-rental-api-gateway/docs"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.Cors())
	s.router.Use(requestid.New())
	s.router.Use(middleware.Base())
	s.router.Use(middleware.Logger(s.log))
}

func (s *Server) setupRoutes(
	tokenManager TokenParser,
	userPermissionsCache UserPermissionsCache,
	userSessionCache UserSessionCache,
) {
	publicV1 := s.router.Group("/api/v1")
	{
		publicV1.GET("/health", s.healthHandler.Health)

		auth := publicV1.Group("/auth")
		{
			auth.POST("/register", s.userHandler.Register)
			auth.POST("/sign-in", s.userHandler.SignIn)
			auth.POST("/refresh-token", s.userHandler.RefreshToken)
		}
	}

	authentication := middleware.NewAuthentication(tokenManager, userPermissionsCache, userSessionCache)

	protectedV1 := s.router.Group("/api/v1")
	protectedV1.Use(authentication.Middleware())
	protectedV1.Use(middleware.SuspensionChecker())
	{
		auth := protectedV1.Group("/auth")
		{
			auth.POST("/sign-out", s.userHandler.SignOut)
		}

		users := protectedV1.Group("/users")
		{
			users.POST("", s.userHandler.Create)
			users.GET("/:id", s.userHandler.Get)
			users.GET("", s.userHandler.GetAllWithFilter)
			users.PATCH("/:id", s.userHandler.Update)
			users.DELETE("/:id", s.userHandler.Delete)

			users.GET("/me", s.userHandler.Me)

			users.POST("/activation-code/send", s.userHandler.SendActivationCode)
			users.POST("/activation-code/check", s.userHandler.CheckActivationCode)

			users.POST("/documents", s.userHandler.CreateDocument)
			users.POST("/documents/upload", s.userHandler.GetUploadDocumentData)
			users.GET("/:id/documents/processed", s.userHandler.GetProcessedDocumentsForUser)
			users.POST("/documents/check/:id", s.userHandler.CheckDocument)
		}

		carModels := protectedV1.Group("/car-models")
		{
			carModels.POST("", s.carModelHandler.Create)
			carModels.GET("/:id", s.carModelHandler.Get)
			carModels.GET("", s.carModelHandler.GetAll)
			carModels.PATCH("/:id", s.carModelHandler.Update)
			carModels.DELETE("/:id", s.carModelHandler.Delete)
			carModels.GET("/image-upload", s.carModelHandler.GetImageUploadUrl)
		}

		verified := protectedV1.Group("")
		verified.Use(middleware.EmailVerificationChecker())
		verified.Use(middleware.DocumentVerificationChecker())
		{
			cars := verified.Group("/cars")
			{
				cars.POST("", s.carHandler.Create)
				cars.GET("/:id", s.carHandler.Get)
				cars.GET("", s.carHandler.GetAll)
				cars.PATCH("/:id", s.carHandler.Update)
				cars.DELETE("/:id", s.carHandler.Delete)
				cars.GET("/status-log", s.carHandler.GetCarStatusLog)
				cars.GET("/fuel-history", s.carHandler.GetCarFuelHistory)
				cars.GET("/image-upload", s.carHandler.GetImageUploadUrl)
			}

			carInsurances := verified.Group("/car-insurances")
			{
				carInsurances.POST("", s.carInsuranceHandler.Create)
				carInsurances.GET("/:id", s.carInsuranceHandler.Get)
				carInsurances.GET("", s.carInsuranceHandler.GetAll)
				carInsurances.PATCH("/:id", s.carInsuranceHandler.Update)
				carInsurances.DELETE("/:id", s.carInsuranceHandler.Delete)
				carInsurances.GET("/image-upload", s.carInsuranceHandler.GetImageUploadUrl)
			}

			carMaintenance := verified.Group("/car-maintenance")
			{
				template := carMaintenance.Group("/template")
				{
					template.POST("", s.carMaintenanceHandler.CreateTemplate)
					template.GET("/:id", s.carMaintenanceHandler.GetTemplate)
					template.GET("", s.carMaintenanceHandler.GetAllTemplates)
					template.PATCH("/:id", s.carMaintenanceHandler.UpdateTemplate)
					template.DELETE("/:id", s.carMaintenanceHandler.DeleteTemplate)
				}

				records := carMaintenance.Group("/records")
				{
					records.GET("", s.carMaintenanceHandler.GetRecords)
					records.POST("/complete/:id", s.carMaintenanceHandler.CompleteRecord)
					records.GET("/receipt-image-upload", s.carMaintenanceHandler.GetReceiptImageUploadUrl)
				}
			}

			zones := verified.Group("/zones")
			{
				zones.POST("", s.zoneHandler.Create)
				zones.GET("/:id", s.zoneHandler.Get)
				zones.GET("", s.zoneHandler.GetAll)
				zones.PATCH("/:id", s.zoneHandler.Update)
				zones.DELETE("/:id", s.zoneHandler.Delete)
			}
		}
	}

	// Swagger
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
