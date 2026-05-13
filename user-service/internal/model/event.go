package model

type DocumentAnalyzedEvent struct {
	DocumentID string
	Passed     bool
	Defects    []Defect
}

type Defect struct {
	Type        string
	Description string
}
