package delivery

type sendFailure struct {
	err error
}

func (f sendFailure) Error() string {
	return f.err.Error()
}

func (f sendFailure) Unwrap() error {
	return f.err
}
