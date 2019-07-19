package types

// Query Result Payload for a get query
type QueryResGet struct {
	Author string `json:"author"`
	Url    string `json:"url"`
}

// implement fmt.Stringer
func (r QueryResGet) String() string {
	return r.Url
}
