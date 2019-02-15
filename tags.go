package tracefall

type Tags []string

func (t Tags) List() []string {
	return []string(t)
}

func (t *Tags) Add(tag string) *Tags {
	*t = append(*t, tag)
	return t
}

func (t *Tags) Clear() *Tags {
	*t = []string{}
	return t
}
