package error

type ErrNoData struct {
	message string
}

func NewErrNoData(msg string) *ErrNoData {
	return &ErrNoData{msg}
}

func (e *ErrNoData) Error() string {
	return e.message
}

type ErrConversion struct {
	message string
}

func NewErrConversion(msg string) *ErrConversion {
	return &ErrConversion{msg}
}

func (e *ErrConversion) Error() string {
	return e.message
}

type ErrDuplicateData struct {
	message string
}

func NewErrDuplicateData(msg string) *ErrDuplicateData {
	return &ErrDuplicateData{msg}
}

func (e *ErrDuplicateData) Error() string {
	return e.message
}
