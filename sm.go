package sm

type Machineable interface {
	GetState() string
	SetState(string)
}

type transitionWhen int

const (
	tunknown = iota
	tbefore
	tafter
	taround
)

type actionFunc func(Machineable) error
type transition struct {
	from   string
	to     string
	action actionFunc
	when   transitionWhen
}

type StateMachine struct {
	states       []string
	transitions  map[string]map[string][]actionFunc
	itransitions map[string]map[string][]actionFunc
}

func New() *StateMachine {
	return &StateMachine{
		states:       []string{},
		transitions:  map[string]map[string][]actionFunc{},
		itransitions: map[string]map[string][]actionFunc{},
	}
}

func (s *StateMachine) Before(t *transition) {
	t.when = tbefore
	addToStack(s.transitions, t.from, t.to, t.action)
	addToStack(s.itransitions, t.to, t.from, t.action)
}

func (s *StateMachine) After(t *transition) {
	t.when = tafter
	addToStack(s.transitions, t.from, t.to, t.action)
	addToStack(s.itransitions, t.to, t.from, t.action)
}

func FromTo(from, to string, f actionFunc) *transition {
	return &transition{
		from:   from,
		to:     to,
		action: f,
	}
}

func FromAny(from string, f actionFunc) *transition {
	return &transition{
		from:   from,
		action: f,
	}
}

func AnyTo(to string, f actionFunc) *transition {
	return &transition{
		to:     to,
		action: f,
	}
}

func (s *StateMachine) Change(obj Machineable, newState string) error {
	oldState := obj.GetState()
	funcs := extractFromStack(s.transitions, oldState, newState)

	return nil
}

func addToStack(stack map[string]map[string][]actionFunc, a, b string, f actionFunc) {
	fl, ok := stack[a]
	if !ok {
		fl = map[string][]actionFunc{}
	}
	sl, ok := fl[b]
	if !ok {
		sl = []actionFunc{}
	}
	sl = append(sl, f)
	fl[b] = sl
	stack[a] = fl
}

func extractFromStack(stack map[string]map[string][]actionFunc, a, b string) []actionFunc {
	funcs := []actionFunc{}
	froms, ok := stack[a]
	if ok {
		straight, ok := froms[b]
		if ok {
			funcs = append(funcs, straight...)
		}
		any, ok := froms[""]
		if ok {
			funcs = append(funcs, any...)
		}
	}
	froms, ok = stack[""]
	if ok {
		any, ok := froms[b]
		if ok {
			funcs = append(funcs, any...)
		}
	}
	return funcs
}
