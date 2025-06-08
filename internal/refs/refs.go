package refs

import "strings"

func ParseRef(ref string) (def, name string, ok bool) {
	elems := strings.Split(ref, "/")
	if len(elems) != 4 {
		return "", "", false
	}
	def, name = elems[2], elems[3]
	return def, name, true
}
