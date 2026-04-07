package http

import (
	_ "github.com/sorawaslocked/car-rental-api-gateway/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/refresh-token", s.authHandler.RefreshToken)
			auth.POST("/logout", s.authHandler.Logout)
		}

		users := v1.Group("/users")
		{
			users.POST("", s.userHandler.Create)
			users.GET("", s.userHandler.Get)
			users.PATCH("", s.userHandler.Update)
			users.DELETE("", s.userHandler.Delete)
			users.GET("/me", s.userHandler.Me)
			users.POST("/activation-code/send", s.userHandler.SendActivationCode)
			users.POST("/activation-code/check", s.userHandler.CheckActivationCode)
		}

		carModels := v1.Group("/car-models")
		{
			carModels.POST("", s.carModelHandler.Create)
			carModels.GET("/:id", s.carModelHandler.Get)
			carModels.GET("", s.carModelHandler.GetAll)
			carModels.PATCH("/:id", s.carModelHandler.Update)
			carModels.DELETE("/:id", s.carModelHandler.Delete)
			carModels.GET("/image-upload", s.carModelHandler.GetImageUploadUrl)
		}

		cars := v1.Group("/cars")
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

		carInsurances := v1.Group("/car-insurances")
		{
			carInsurances.POST("", s.carInsuranceHandler.Create)
			carInsurances.GET("/:id", s.carInsuranceHandler.Get)
			carInsurances.GET("", s.carInsuranceHandler.GetAll)
			carInsurances.PATCH("/:id", s.carInsuranceHandler.Update)
			carInsurances.DELETE("/:id", s.carInsuranceHandler.Delete)
			carInsurances.GET("/image-upload", s.carInsuranceHandler.GetImageUploadUrl)
		}

		carMaintenance := v1.Group("/car-maintenance")
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

		zones := v1.Group("/zones")
		{
			zones.POST("", s.zoneHandler.Create)
			zones.GET("/:id", s.zoneHandler.Get)
			zones.GET("", s.zoneHandler.GetAll)
			zones.PATCH("/:id", s.zoneHandler.Update)
			zones.DELETE("/:id", s.zoneHandler.Delete)
		}
	}

	// Swagger
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
