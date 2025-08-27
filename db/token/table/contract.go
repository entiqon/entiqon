package table

import "github.com/entiqon/entiqon/db/contract"

type Token interface {
	contract.BaseToken
	contract.Clonable[Token]
	contract.Debuggable
	contract.Errorable[Token]
	contract.Rawable
	contract.Renderable
	contract.Stringable
	contract.Validable
	Name() string
}

var _ Token = (*Table)(nil)
