package loli2

// Noise is noise instance in lolinote.
type Noise struct {
	path string
}

// Path return noise's path.
func (n Noise) Path() string {
	return n.path
}
