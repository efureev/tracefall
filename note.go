package tracefall

import (
	"encoding/json"
	"time"
)

// Note struct
type Note struct {
	Time int64  `json:"t"`
	Note string `json:"v"`
}

// NewNote struct
func NewNote(note string) *Note {
	return &Note{time.Now().UnixNano(), note}
}

// Notes list of Note
type Notes []*Note

// NoteGroup struct
type NoteGroup struct {
	Notes Notes  `json:"notes"`
	Label string `json:"label"`
}

// Count elements in NoteGroup
func (n NoteGroup) Count() int {
	return len(n.Notes)
}

// Add note to NoteGroup
func (n *NoteGroup) Add(note string) *NoteGroup {
	n.Notes = append(n.Notes, NewNote(note))
	return n
}

// NewNoteGroup create new NoteGroup
func NewNoteGroup(groupLabel string) *NoteGroup {
	return &NoteGroup{Label: groupLabel}
}

// Clear NoteGroup from notes
func (n *NoteGroup) Clear() *NoteGroup {
	n.Notes = Notes{}
	return n
}

// NoteGroupList is list of NoteGroup
type NoteGroupList []*NoteGroup

// FromJSON fill NoteGroupList from json
func (n *NoteGroupList) FromJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// NoteGroups is list of NoteGroup
type NoteGroups map[string]*NoteGroup

// Count elements in NoteGroups
func (n NoteGroups) Count() int {
	return len(n)
}

// NewNotesGroups creates new NoteGroups struct
func NewNotesGroups() NoteGroups {
	return make(NoteGroups)
}

// Add new note to exist group (or create new if absent) in NoteGroups
func (n NoteGroups) Add(group, note string) NoteGroups {
	if lg, ok := n[group]; ok {
		n[group] = lg.Add(note)
	} else {
		n[group] = NewNoteGroup(group).Add(note)
	}
	return n
}

// AddGroup add notes list to exist group (or create new if absent) in NoteGroups
func (n NoteGroups) AddGroup(group string, notes []string) NoteGroups {
	for _, note := range notes {
		n.Add(group, note)
	}

	return n
}

// AddNoteGroup add group struct to list
func (n NoteGroups) AddNoteGroup(group *NoteGroup) NoteGroups {
	if group != nil {
		n[group.Label] = group
	}

	return n
}

// Get NoteGroup from NoteGroup list
func (n NoteGroups) Get(groupName string) *NoteGroup {
	if lg, ok := n[groupName]; ok {
		return lg
	}
	return nil
}

// Remove NoteGroup from list
func (n NoteGroups) Remove(groupName string) NoteGroups {
	delete(n, groupName)
	return n
}

// Clear NoteGroups list
func (n *NoteGroups) Clear() *NoteGroups {
	*n = NewNotesGroups()
	return n
}

// ToJSON return json bytes of NoteGroups
func (n NoteGroups) ToJSON() []byte {
	b, err := json.Marshal(n.prepareToJSON())
	if err != nil {
		b = []byte(`{}`)
	}
	return b
}

func (n NoteGroups) prepareToJSON() NoteGroupList {
	var list NoteGroupList
	for _, ng := range n {
		list = append(list, ng)
	}
	return list
}

// ToJSONString return json string of NoteGroups
func (n NoteGroups) ToJSONString() string {
	return string(n.ToJSON())
}

// FromJSON fill NoteGroups from json. Previously clear list
func (n *NoteGroups) FromJSON(b []byte) error {
	n.Clear()

	var list NoteGroupList
	err := json.Unmarshal(b, &list)
	if err != nil {
		return err
	}

	for _, ng := range list {
		n.AddNoteGroup(ng)
	}

	return nil
}
