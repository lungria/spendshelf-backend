package account

type Type int

const (
	Regular Type = iota
	Savings
)

func (d Type) String() string {
	return [...]string{"Regular", "Savings"}[d]
}
