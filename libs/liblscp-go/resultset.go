package liblscp

type ResultSet struct {
	Index           int
	Code            int
	Message         string
	IsWarning       bool
	Result          string
	MultiLineResult []string
}

func (r *ResultSet) AddLine(ln string) {
	r.MultiLineResult = append(r.MultiLineResult, ln)
}
