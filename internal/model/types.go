package model

type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Path      string
	Labels    map[string]string
	Selector  map[string]string
}
