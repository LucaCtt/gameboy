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
