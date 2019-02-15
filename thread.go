package tracefall

// Thread is list of LogJSON struct
type Thread []*LogJSON

// Add Log to Thread
func (t *Thread) Add(log *LogJSON) Thread {
	*t = append(*t, log)
	return *t
}

// ThreadFromList make Thread from list of LogJSON structs
func ThreadFromList(list []*LogJSON) Thread {
	return Thread(list)
}
