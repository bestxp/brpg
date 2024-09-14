package actions

type Action string

const (
	UnitIdle Action = "idle"
	UnitMove Action = "run"
)

func (a Action) String() string {
	return string(a)
}
