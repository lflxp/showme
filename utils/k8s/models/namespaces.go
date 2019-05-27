package models

type Namespace struct {
	Name string
}

// mutil Names for delete
type Namespaces struct {
	Names []string `json:"name"`
}
