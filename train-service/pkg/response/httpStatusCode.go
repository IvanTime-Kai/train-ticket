package response

const (
	ErrCodeSuccess       = 20000
	ErrCodeParamsInvalid = 40000
	ErrEmailExists       = 40001
	ErrInvalidCreds      = 40002
	ErrNotFound          = 40004
	ErrUnauthorized      = 40003
	ErrInternalServer    = 50000
)

var msg = map[int]string{
	ErrCodeSuccess:       "Success",
	ErrCodeParamsInvalid: "Params invalid",
	ErrEmailExists:       "Email already exists",
	ErrInvalidCreds:      "Invalid email or password",
	ErrNotFound:      "User not found",
	ErrUnauthorized:      "Unauthorized",
	ErrInternalServer:    "Internal Server",
}
