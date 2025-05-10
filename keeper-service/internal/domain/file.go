package domain

type File struct {
	Data     string
	Name     string
	Hash     string
	Location string
}

func NewFile(data, name, hash, location string) *File {
	return &File{
		Data:     data,
		Name:     name,
		Hash:     hash,
		Location: location,
	}
}
