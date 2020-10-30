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
	CodeUnrecognized                 ErrorCode = "M_UNRECOGNIZED"
	CodeUnauthorized                 ErrorCode = "M_UNAUTHORIZED"
	CodeUserDeactivated              ErrorCode = "M_USER_DEACTIVATED"
	CodeUserInUse                    ErrorCode = "M_USER_IN_USE"
	CodeInvalidUsername              ErrorCode = "M_INVALID_USERNAME"
	CodeRoomInUse                    ErrorCode = "M_ROOM_IN_USE"
	CodeInvalidRoomState             ErrorCode = "M_INVALID_ROOM_STATE"
	CodeThreePIDInUse                ErrorCode = "M_THREEPID_IN_USE"
	CodeThreePIDNotFound             ErrorCode = "M_THREEPID_NOT_FOUND"
	CodeThreePIDAuthFailed           ErrorCode = "M_THREEPID_AUTH_FAILED"
	CodeThreePIDDenied               ErrorCode = "M_THREEPID_DENIED"
	CodeServerNotTrusted             ErrorCode = "M_SERVER_NOT_TRUSTED"
	CodeUnsupportedRoomVersion       ErrorCode = "M_UNSUPPORTED_ROOM_VERSION"
	CodeIncompatibleRoomVersion      ErrorCode = "M_INCOMPATIBLE_ROOM_VERSION"
	CodeBadState                     ErrorCode = "M_BAD_STATE"
	CodeGuestAccessForbidden         ErrorCode = "M_GUEST_ACCESS_FORBIDDEN"
	CodeCaptchaNeeded                ErrorCode = "M_CAPTCHA_NEEDED"
	CodeMissingParam                 ErrorCode = "M_MISSING_PARAM"
	CodeInvalidParam                 ErrorCode = "M_INVALID_PARAM"
	CodeTooLarge                     ErrorCode = "M_TOO_LARGE"
	CodeExclusive                    ErrorCode = "M_EXCLUSIVE"
	CodeResourceLimitExceeded        ErrorCode = "M_RESOURCE_LIMIT_EXCEEDED"
	CodeCannotLeaveServiceNoticeRoom ErrorCode = "M_CANNOT_LEAVE_SERVICE_NOTICE_ROOM"

	// Codes that are documented on other sections
	CodeWeakPassword ErrorCode = "M_WEAK_PASSWORD"
)
