package auth

const (
	roleAdmin = "admin"
	roleUser  = "user"
	roleGuest = "guest"
)

const (
	carModelServiceCreateCarModel = "/service.car.CarModelService/CreateCarModel"
	carModelServiceGetCarModel    = "/service.car.CarModelService/GetCarModel"
	carModelServiceGetCarModels   = "/service.car.CarModelService/GetCarModels"
	carModelServiceUpdateCarModel = "/service.car.CarModelService/UpdateCarModel"
	carModelServiceDeleteCarModel = "/service.car.CarModelService/DeleteCarModel"

	carServiceCreateCar        = "/service.car.CarService/CreateCar"
	carServiceGetCar           = "/service.car.CarService/GetCar"
	carServiceGetCars          = "/service.car.CarService/GetCars"
	carServiceGetAvailableCars = "/service.car.CarService/GetAvailableCars"
	carServiceUpdateCar        = "/service.car.CarService/UpdateCar"
	carServiceUpdateCarStatus  = "/service.car.CarService/UpdateCarStatus"
	carServiceDeleteCar        = "/service.car.CarService/DeleteCar"
)

var methodPermissions = map[string][]string{
	carModelServiceCreateCarModel: {roleAdmin},
	carModelServiceUpdateCarModel: {roleAdmin},
	carModelServiceDeleteCarModel: {roleAdmin},
	carModelServiceGetCarModel:    {roleAdmin, roleUser, roleGuest},
	carModelServiceGetCarModels:   {roleAdmin, roleUser, roleGuest},

	carServiceCreateCar:        {roleAdmin},
	carServiceUpdateCar:        {roleAdmin},
	carServiceUpdateCarStatus:  {roleAdmin},
	carServiceDeleteCar:        {roleAdmin},
	carServiceGetCar:           {roleAdmin},
	carServiceGetCars:          {roleAdmin},
	carServiceGetAvailableCars: {roleAdmin, roleUser, roleGuest},
}

func isAllowed(method string, roles []string) bool {
	allowedRoles, ok := methodPermissions[method]
	if !ok {
		return false
	}

	allowedSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowedSet[role] = struct{}{}
	}

	for _, role := range roles {
		if _, ok := allowedSet[role]; ok {
			return true
		}
	}

	return false
}
