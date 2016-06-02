package validation

type Errs struct {
	verrs   map[string][]string
	hasErrs bool
}

func Errors() Errs {
	return Errs{
		verrs:   make(map[string][]string),
		hasErrs: false,
	}
}

func (e Errs) Add(field, errStr string) {
	e.hasErrs = true
	if _, ok := e.verrs[field]; !ok {
		e.verrs[field] = make([]string, 0)
	}
	e.verrs[field] = append(e.verrs[field], errStr)
}

func (e Errs) HasErrs() bool { return e.hasErrs }
