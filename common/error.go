package common

type ignore struct {
	err error
}

// https://godoc.org/github.com/pkg/errors
type wrapped interface {
	Cause() error
}

func (e ignore) Cause() error {
	return e.err
}

func (e ignore) Error() string {
	return e.err.Error()
}

func MakeIgnore() error {
	return ignore{}
}

func MakeIgnoreErr(err error) error {
	return ignore{err: err}
}

func TraceBack(err error) error {
	for e := err; e != nil; {
		switch e.(type) {
		case ignore:
			return nil
		case wrapped:
			e = e.(wrapped).Cause()
		default:
			return e
		}
	}

	return nil
}
