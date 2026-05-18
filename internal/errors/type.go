package errors

type AppError struct {
	Code    int
	Message string
	Actual  error
}

func (err AppError) Error() string {
	if err.Actual != nil {
		return err.Message + ": " + err.Actual.Error()
	}
	return err.Message

}

func (err AppError) Wrap(actual error) AppError {
	err.Actual = actual
	return err
}

func (err AppError) Unwrap() error {
	return err.Actual
}
