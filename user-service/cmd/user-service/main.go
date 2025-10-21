package main

import (
	"car-rental-user-service/internal/app"
	"car-rental-user-service/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg)
	fmt.Println(application)
}
