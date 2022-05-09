package register

type Register interface {
	Register() (err error)
	Listen()
	Close() (err error)
}
