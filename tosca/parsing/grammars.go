package parsing

//
// Grammar
//

type Grammar struct {
	Readers                    Readers
	InvalidNamespaceCharacters string
}

func NewGrammar() Grammar {
	return Grammar{
		Readers: make(Readers),
	}
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
