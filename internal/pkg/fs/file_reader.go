package fs

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FileReaderStub struct {
	bytes []byte
	err   error
}

func NewFileReaderStub(bytes []byte, err error) *FileReaderStub {
	return &FileReaderStub{
		bytes: bytes,
		err:   err,
	}
}

func (r *FileReaderStub) ReadFile(string) ([]byte, error) {
	return r.bytes, r.err
}
