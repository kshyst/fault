package ftag

const (
	None             Kind = ""                  // Empty error.
	Internal         Kind = "INTERNAL"          // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
	Cancelled        Kind = "CANCELLED"         // The operation was cancelled, typically by the caller.
	InvalidArgument  Kind = "INVALID_ARGUMENT"  // The client specified an invalid argument.
	NotFound         Kind = "NOT_FOUND"         // Some requested entity was not found.
	AlreadyExists    Kind = "ALREADY_EXISTS"    // The entity that a client attempted to create already exists.
	PermissionDenied Kind = "PERMISSION_DENIED" // The caller does not have permission to execute the specified operation.
	Unauthenticated  Kind = "UNAUTHENTICATED"   // The request does not have valid authentication credentials for the operation.

	ValidationError Kind = "VALIDATION_ERROR" // The request was invalid, but not in a way that is covered by the other error codes. This is used for validation errors that are not covered by the other error codes.
)
