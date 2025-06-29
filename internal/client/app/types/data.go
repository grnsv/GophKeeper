package types

type Credentials struct {
	Resource string
	Login    string
	Password string
	Metadata map[string]string
}

type Text struct {
	Text     string
	Metadata map[string]string
}

type Binary struct {
	Binary   []byte
	Metadata map[string]string
}

type Card struct {
	CardNumber string
	EXP        string
	CVV        string
	CardHolder string
	Metadata   map[string]string
}
