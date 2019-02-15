package tracefall

// Tags struct is list of string
type Tags []string

// List of tags
func (t Tags) List() []string {
	return []string(t)
}

// Add new Tag by name
func (t *Tags) Add(tag string) *Tags {
	*t = append(*t, tag)
	return t
}

// Clear tag list
func (t *Tags) Clear() *Tags {
	*t = []string{}
	return t
}
