package traceFall

type Thread []*LogJSON

func (t *Thread) Add(log *LogJSON) Thread {
	*t = append(*t, log)
	return *t
}

func ThreadFromList(list []*LogJSON) Thread {
	return Thread(list)
}
