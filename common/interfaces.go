package common

type IErrorUtils interface {
	Error(err error)
	ErrorStr(err string)
	FatalError(err error)
	FatalErrorStr(err string)
}