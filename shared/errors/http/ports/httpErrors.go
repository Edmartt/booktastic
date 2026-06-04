package ports

type HTTPErrors interface {
	WriteError(code int, message string)
}
