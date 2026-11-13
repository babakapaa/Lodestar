package league

type Display string

const (
	DisplayName     Display = "name"
	DisplayHidden   Display = "hidden"
	DisplayNickname Display = "nickname"
)

func (d Display) valid() bool {
	return d == DisplayName || d == DisplayHidden || d == DisplayNickname
}
