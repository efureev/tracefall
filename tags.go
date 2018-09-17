package traceFall

type Tags []string

func (t Tags) List() []string {
	return []string(t)
}

func (t *Tags) Add(tag string) {
	*t = append(*t, tag)
}

func (t *Tags) Clear() {
	*t = []string{}
}
