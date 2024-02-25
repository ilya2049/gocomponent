package fs

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}
