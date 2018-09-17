package traceFall

import "github.com/satori/go.uuid"

func generateUUID() uuid.UUID {
	if uid, err := uuid.NewV4(); err != nil {
		return generateUUID()
	} else {
		return uid
	}
}

func removeDuplicatesFromSlice(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}
