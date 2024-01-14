package error

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GLCResponseError model for returning errors.
type GLCResponseError struct {
	HTTPStatusCode     int       `json:"httpStatusCode"`
	ErrorCode          ErrorCode `json:"errorCode"`
	Msg                string    `json:"message"`
	RecommendedActions []string  `json:"recommendedActions"`
}

// errInfoMap represents mapping of errorResponse struct with errorCode.
// nolint:gochecknoglobals, unused
var errInfoMap = map[ErrorCode]*GLCResponseError{
	FailedUnmarshalling: {
		HTTPStatusCode:     http.StatusBadRequest,
		ErrorCode:          FailedUnmarshalling,
		Msg:                "Failed unmarshalling payload content",
		RecommendedActions: []string{"Reverify json payload"},
	},

	FailedMarshaling: {
		HTTPStatusCode:     http.StatusBadRequest,
		ErrorCode:          FailedMarshaling,
		Msg:                "Failed marshaling data",
		RecommendedActions: []string{"Reverify the provided data"},
	},

	DataParsingFailed: {
		HTTPStatusCode:     http.StatusInternalServerError,
		ErrorCode:          DataParsingFailed,
		Msg:                "Failed parsing data values",
		RecommendedActions: []string{"Reverify the provided data"},
	},

	MissingXTenantID: {
		HTTPStatusCode:     http.StatusBadRequest,
		ErrorCode:          MissingXTenantID,
		Msg:                "X-Tenant-ID is missing",
		RecommendedActions: []string{"Pass the X-Tenant-ID while making request"},
	},

	Forbidden: {
		HTTPStatusCode:     http.StatusForbidden,
		ErrorCode:          Forbidden,
		Msg:                "User not authorized to perform operation",
		RecommendedActions: []string{"Check if user has the required permissions"},
	},

	FailedDataValidation: {
		HTTPStatusCode:     http.StatusBadRequest,
		ErrorCode:          FailedDataValidation,
		Msg:                "Failed validating the data content",
		RecommendedActions: []string{"Reverify the provided data"},
	},

	MethodNotAllowed: {
		HTTPStatusCode:     http.StatusMethodNotAllowed,
		ErrorCode:          MethodNotAllowed,
		Msg:                "Invalid method type requested",
		RecommendedActions: []string{"Method type not supported"},
	},

	BadRequest: {
		HTTPStatusCode:     http.StatusBadRequest,
		ErrorCode:          BadRequest,
		Msg:                "Server cannot process the request",
		RecommendedActions: []string{"Reverify the provided request"},
	},
}

// nolint:unused
func RespondWithError(ginCtx *gin.Context, err error, msg string) {
	status := http.StatusInternalServerError
	errorCode := InternalServerError
	message := "Internal server error"
	recommendActions := []string{"Contact admin for support"}

	errType := GetErrorType(err)

	reflectErr, ok := errInfoMap[errType]
	if ok {
		status = reflectErr.HTTPStatusCode
		errorCode = reflectErr.ErrorCode
		message = reflectErr.Msg
		recommendActions = reflectErr.RecommendedActions
	}

	if msg != "" {
		message = msg
	}

	ginCtx.JSON(status, gin.H{
		"HTTPStatusCode":     status,
		"errorCode":          errorCode,
		"message":            message,
		"recommendedActions": recommendActions,
	})
}
