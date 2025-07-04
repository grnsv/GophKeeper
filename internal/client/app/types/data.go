package types

type Metadata map[string]string

type Credentials struct {
	Resource string
	Login    string
	Password string
	Metadata Metadata
}

type Text struct {
	Text     string
	Metadata Metadata
}

type Binary struct {
	Binary   []byte
	Metadata Metadata
}

type Card struct {
	CardNumber string
	EXP        string
	CVV        string
	CardHolder string
	Metadata   Metadata
}

type Data interface {
	Credentials | Text | Binary | Card
}
