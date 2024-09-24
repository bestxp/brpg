package scene

import e "github.com/hajimehoshi/ebiten/v2"

type SceneID int

const (
	Welcome SceneID = iota
	Game
)

type GUIEvent int

const (
	Click GUIEvent = iota
	Hover
)

type EventFn func() error

type EventStore map[string]map[GUIEvent][]EventFn

func NewEventStore() EventStore {
	return make(EventStore)
}

func (s EventStore) Get(name string, eventID GUIEvent) []EventFn {
	if _, ok := s[name]; !ok {
		return nil
	}
	return s[name][eventID]
}

func (s EventStore) Add(name string, eventID GUIEvent, fn EventFn) {
	if _, ok := s[name]; !ok {
		s[name] = map[GUIEvent][]EventFn{}
	}
	if _, ok := s[name][eventID]; !ok {
		s[name][eventID] = make([]EventFn, 0)
	}
	s[name][eventID] = append(s[name][eventID], fn)
}

type Scene interface {
	DrawAt(screen *e.Image)
	MouseMove(x, y int)
	Frame(frame int)
	Update()
}
