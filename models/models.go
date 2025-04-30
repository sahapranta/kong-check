package models

type Route struct {
	ID          string
	Name        string
	Paths       []string
	Methods     []string
	Protocols   []string
	ServiceID   string
	ServiceName string
	HostNames   []string
	Headers     map[string][]string
}

type Service struct {
	ID   string
	Name string
	Host string
	Port int
	Path string
}

type CheckResult struct {
	ServiceName  string
	RouteName    string
	URL          string
	Method       string
	Status       int
	Error        error
	ResponseTime int64  // in milliseconds
	ResponseBody string // Only filled when verbose is true
}
