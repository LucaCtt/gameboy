package errors

// E creates a new error with the given args.
// Will panic if an arg not supported by Error is passed.
func E(args ...interface{}) *Error {
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Message = arg
		case ErrCode:
			e.Code = arg
		case ErrComponent:
			e.Component = arg
		case error:
			e.Err = arg
		default:
			panic("invalid arg")
		}
	}

	return e
}

// Code returns the errors code.
// If it does not exist (which means err is not a pointer to Error) CodeUnexpected will be returned.
// If the code exists but is 0, the wrapped error's code will be returned.
func Code(err error) ErrCode {
	e, ok := err.(*Error)
	if !ok {
		return CodeUnexpected
	}

	if e.Code != 0 {
		return e.Code
	}

	return Code(e.Err)
}
