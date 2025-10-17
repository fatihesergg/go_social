package errors

type AppError struct {
	Code    int
	Message string
	Actual  error
}

func (err AppError) Error() string {
	return err.Message
}

func (err AppError) ActualErr() string {
	if err.Actual != nil {
		return err.Actual.Error()
	}
	return ""

}

func (err AppError) Wrap(actual error) AppError {
	err.Actual = actual
	return err
}
