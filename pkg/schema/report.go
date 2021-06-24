package schema

import "time"

type Report struct {
	Id               string `json:"id"`
	Source           string `json:"source"`
	SourceIdentityId string `json:"sourceIdentityId"`
	Reference        struct {
		Id   string `json:"referenceId"`
		Type string `json:"referenceType"`
	}
	State   string `json:"state"`
	Payload struct {
		Source                string `json:"source"`
		ReportType            string `json:"reportType"`
		Message               string `json:"message"`
		ReportId              string `json:"reportId"`
		ReferenceResourceId   string `json:"referenceResourceId"`
		ReferenceResourceType string `json:"referenceResourceType"`
	} `json:"payload"`
	Created time.Time `json:"created"`
}
