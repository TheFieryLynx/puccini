package parsing

type DataDefinitionValidator func(EntityPtr)
type DataValueValidator func(*Context, EntityPtr, EntityPtr)

//
// Grammar
//

type Grammar struct {
	Readers                    Readers
	InvalidNamespaceCharacters string
	DataDefinitionValidator    DataDefinitionValidator
	DataValueValidator         DataValueValidator
}

func NewGrammar() Grammar {
	return Grammar{
		Readers: make(Readers),
	}
}

func (self *Grammar) RegisterReader(name string, reader Reader) {
	self.Readers[name] = reader
}
