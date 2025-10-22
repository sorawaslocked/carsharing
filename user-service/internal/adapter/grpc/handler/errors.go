package handler

func fromErrToStringMap(m map[string]error) map[string]string {
	newMap := make(map[string]string)

	for k, v := range m {
		newMap[k] = v.Error()
	}

	return newMap
}
