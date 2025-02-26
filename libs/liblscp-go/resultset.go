package liblscp

type ResultTypeType int

const (
	RsOk ResultTypeType = iota
	RsError
	RsWarning
)

var ResultType = struct {
	Ok      ResultTypeType
	Error   ResultTypeType
	Warning ResultTypeType
}{
	Ok:      RsOk,
	Error:   RsError,
	Warning: RsWarning,
}

type ResultSet struct {
	Type            ResultTypeType
	Index           int
	Code            int
	IsMultiline     bool
	Message         string
	MultiLineResult []string
}

func (r *ResultSet) AddLine(ln string) {
	r.MultiLineResult = append(r.MultiLineResult, ln)
	r.IsMultiline = true
}
