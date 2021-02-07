package shortener

type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(*Redirect) ([]byte, error)
}