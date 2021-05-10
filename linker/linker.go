package linker

type Linker interface {
	// Attempts to create a link between a provided ressource source and destination
	Link() (int, int, error)
}

func Link(l Linker) (int, int, error) {
	return l.Link()
}
