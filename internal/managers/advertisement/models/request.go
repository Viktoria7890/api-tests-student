package models

// Photo одно изображение, прикладываемое к объявлению при его создании
type Photo struct {
	FileName string
	Content  []byte
}
