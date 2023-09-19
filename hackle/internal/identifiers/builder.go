package identifiers

type Builder struct {
	identifiers map[string]string
}

func NewBuilder() *Builder {
	return &Builder{
		identifiers: make(map[string]string),
	}
}

func (b *Builder) Add(identifierType string, identifierValue string) *Builder {
	if b.isValid(identifierType, identifierValue) {
		b.identifiers[identifierType] = identifierValue
	}
	return b
}

func (b *Builder) AddAll(identifiers map[string]string) *Builder {
	for identifierType, identifierValue := range identifiers {
		b.Add(identifierType, identifierValue)
	}
	return b
}

func (b *Builder) Build() map[string]string {
	identifiers := make(map[string]string)
	for t, v := range b.identifiers {
		identifiers[t] = v
	}
	return identifiers
}

func (b *Builder) isValid(identifierType string, identifierValue string) bool {
	if len(identifierType) == 0 {
		return false
	}
	if len(identifierType) > maxIdentifierTypeLength {
		return false
	}
	if len(identifierValue) == 0 {
		return false
	}
	if len(identifierValue) > maxIdentifierValueLength {
		return false
	}
	return true
}

const (
	maxIdentifierTypeLength  = 128
	maxIdentifierValueLength = 512
)
