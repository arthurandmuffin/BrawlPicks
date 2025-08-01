package http

type Code string

const (
	OK    Code = "OK"
	EHTTP Code = "EHTTP"
	EJSON Code = "EJSON"
	ERES  Code = "ERES"
)

func (c Code) Error() string {
	return string(c)
}
