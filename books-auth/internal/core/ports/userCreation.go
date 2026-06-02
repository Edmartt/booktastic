package ports

type AuthProvider interface {
	SignUp(email, password string) (*string, error)
	Login(email, password string) (*string, error)
}
