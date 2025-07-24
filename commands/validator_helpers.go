package commands

func wrapBackquoteEach(values []string) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = "`" + v + "`"
	}
	return result
}
