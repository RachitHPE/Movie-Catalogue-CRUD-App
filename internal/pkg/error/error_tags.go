package error

const (
	// InternalError ...
	InternalError ErrorCode = "Internal Error"

	// FailedUnmarshalling provides error code for unmarshalling failures.
	FailedUnmarshalling ErrorCode = "UNMARSHALLING_FAILED"

	// FailedMarshaling provides error code for marshaling failures.
	FailedMarshaling ErrorCode = "MARSHALING_FAILED"

	// DataParsingFailed provides error code for data parsing failures.
	DataParsingFailed ErrorCode = "FAILED_PARSING_DATA"

	// MissingXTenantID provides error code for missing XTenantID.
	MissingXTenantID ErrorCode = "MISSING_X_Tenant_ID"

	// Forbidden provides error code for unauthorized forbidden calls.
	Forbidden ErrorCode = "FORBIDDEN"

	// FailedDataValidation provides error code for validation failures.
	FailedDataValidation ErrorCode = "DATA_VALIDATION_FAILED"

	// MethodNotAllowed provides error code for method not allowed on a resource.
	MethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"

	// BadRequest provides error code for cases where server cannot process the request.
	BadRequest ErrorCode = "BAD_REQUEST"

	// InternalServerError provides error code for some internal error.
	InternalServerError ErrorCode = "HPE_GL_MP_INTERNAL_ERROR"
)