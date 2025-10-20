package grpc

func grpcError(err error) map[string]string {
	errors := make(map[string]string)
	errors["grpc"] = err.Error()

	return errors
}
