package enum

var enums []PGEnum

type PGEnum struct {
	Name   string
	Values []string
}

func AddEnum(enum PGEnum) {
	enums = append(enums, enum)
}

func GetEnumList() []PGEnum {
	return enums
}