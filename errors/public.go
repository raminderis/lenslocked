package errors

func Publicmyerror(err error, message string) error {
	return publicError{
		err: err,
		msg: message,
	}
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
