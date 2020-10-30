package matrix

// ErrorCode represents an error code that is found in REST errors.
type ErrorCode string

// List of official error codes.
// It can be found at https://matrix.org/docs/spec/client_server/r0.6.1#api-standards.
const (
	// Common error codes
	CodeForbidden     ErrorCode = "M_FORBIDDEN"
	CodeUnknownToken  ErrorCode = "M_UNKNOWN_TOKEN"
	CodeMissingToken  ErrorCode = "M_MISSING_TOKEN"
	CodeBadJSON       ErrorCode = "M_BAD_JSON"
	CodeNotJSON       ErrorCode = "M_NOT_JSON"
	CodeNotFound      ErrorCode = "M_NOT_FOUND"
	CodeLimitExceeded ErrorCode = "M_LIMIT_EXCEEDED"
	CodeUnknown       ErrorCode = "M_UNKNOWN"

	// Other error codes the client might encounter
	CodeUnrecognized                 = "M_UNRECOGNIZED"
	CodeUnauthorized                 = "M_UNAUTHORIZED"
	CodeUserDeactivated              = "M_USER_DEACTIVATED"
	CodeUserInUse                    = "M_USER_IN_USE"
	CodeInvalidUsername              = "M_INVALID_USERNAME"
	CodeRoomInUse                    = "M_ROOM_IN_USE"
	CodeInvalidRoomState             = "M_INVALID_ROOM_STATE"
	CodeThreePIDInUse                = "M_THREEPID_IN_USE"
	CodeThreePIDNotFound             = "M_THREEPID_NOT_FOUND"
	CodeThreePIDAuthFailed           = "M_THREEPID_AUTH_FAILED"
	CodeThreePIDDenied               = "M_THREEPID_DENIED"
	CodeServerNotTrusted             = "M_SERVER_NOT_TRUSTED"
	CodeUnsupportedRoomVersion       = "M_UNSUPPORTED_ROOM_VERSION"
	CodeIncompatibleRoomVersion      = "M_INCOMPATIBLE_ROOM_VERSION"
	CodeBadState                     = "M_BAD_STATE"
	CodeGuestAccessForbidden         = "M_GUEST_ACCESS_FORBIDDEN"
	CodeCaptchaNeeded                = "M_CAPTCHA_NEEDED"
	CodeMissingParam                 = "M_MISSING_PARAM"
	CodeInvalidParam                 = "M_INVALID_PARAM"
	CodeTooLarge                     = "M_TOO_LARGE"
	CodeExclusive                    = "M_EXCLUSIVE"
	CodeResourceLimitExceeded        = "M_RESOURCE_LIMIT_EXCEEDED"
	CodeCannotLeaveServiceNoticeRoom = "M_CANNOT_LEAVE_SERVICE_NOTICE_ROOM"
)
