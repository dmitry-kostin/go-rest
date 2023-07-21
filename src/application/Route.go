package application

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler HandlerFunc
}
