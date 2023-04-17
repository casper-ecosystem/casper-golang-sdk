package sdk

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewRpcError(message string, code int) error {
	return &RpcError{Message: message, Code: code}
}
func (e *RpcError) Error() string {
	return e.Message
}
