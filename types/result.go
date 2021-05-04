package types

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
