package http

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
	}
}
