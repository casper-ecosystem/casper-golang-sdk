package types

// CLValueResult representing a result of an operation that might have failed
type CLValueResult struct {
	IsSuccess bool
	Success   *CLValue
	Error     *CLValue
}

func (r CLValueResult) ResultFieldName() string {
	return "IsSuccess"
}

func (r CLValueResult) SuccessFieldName() string {
	return "Success"
}

func (r CLValueResult) ErrorFieldName() string {
	return "Error"
}
