package traceFall

import (
	"encoding/json"
	"time"
)

type Note struct {
	Time int64  `json:"t"`
	Note string `json:"v"`
}

func NewNote(note string) *Note {
	return &Note{time.Now().UnixNano(), note}
}

type Notes []*Note

type NoteGroup struct {
	Notes Notes  `json:"notes"`
	Label string `json:"label"`
}

func (n NoteGroup) Count() int {
	return len(n.Notes)
}

func (n *NoteGroup) Add(note string) *NoteGroup {
	n.Notes = append(n.Notes, NewNote(note))
	return n
}

func NewNoteGroup(groupLabel string) *NoteGroup {
	return &NoteGroup{Label: groupLabel}
}

func (n *NoteGroup) Clear() *NoteGroup {
	n.Notes = Notes{}
	return n
}

type NoteGroupList []*NoteGroup

func (n *NoteGroupList) FromJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

type NoteGroups map[string]*NoteGroup

func (n NoteGroups) Count() int {
	return len(n)
}

func NewNotesGroups() NoteGroups {
	return make(NoteGroups)
}

func (n NoteGroups) Add(group, note string) NoteGroups {
	if lg, ok := n[group]; ok {
		n[group] = lg.Add(note)
	} else {
		n[group] = NewNoteGroup(group).Add(note)
	}
	return n
}

func (n NoteGroups) AddGroup(group string, notes []string) NoteGroups {
	for _, note := range notes {
		n.Add(group, note)
	}

	return n
}

func (n NoteGroups) Get(groupName string) *NoteGroup {
	if lg, ok := n[groupName]; ok {
		return lg
	}
	return nil
}

func (n NoteGroups) Remove(groupName string) NoteGroups {
	delete(n, groupName)
	return n
}

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

func (n NoteGroups) ToJSONString() string {
	return string(n.ToJSON())
}

func (n *NoteGroups) FromJSON(b []byte) error {
	return json.Unmarshal(b, n)
}
