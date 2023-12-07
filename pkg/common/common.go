package common

type IntegrityController interface {
	AddFileToIc(s string, options ...ICOption) error
	InitDatabase() error
	RefreshIntegrityDatabase() error
	HasIntegrityViolation(filepath string, checkOpts ...string) (bool, map[string]interface{}, error)
}

type ICOption interface {
	Name() string
	Value() interface{}
}
