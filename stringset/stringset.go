package stringset

func IndexOf(list []string, item string) int {
	for i, s := range list {
		if s == item {
			return i
		}
	}
	return -1
}
func Add(list []string, item string) []string {
	i := IndexOf(list, item)
	if i >= 0 {
		return list
	}
	return append(list, item)
}
func Remove(list []string, item string) []string {
	i := IndexOf(list, item)
	if i < 0 {
		return list
	}
	return append(list[:i], list[i+1:]...)
}
