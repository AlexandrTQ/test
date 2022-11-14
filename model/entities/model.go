package entities

var models = []interface{}{
	(*User)(nil),
	(*Transaction)(nil),
}

func GetModels() []interface{} {
	return models
}
