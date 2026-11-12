package module

type Permissioned interface {
	Permissions() []string
}
