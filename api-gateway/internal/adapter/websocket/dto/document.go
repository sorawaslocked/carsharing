package wsdto

type DocumentDefect struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type DocumentAnalyzedMessage struct {
	DocumentID string           `json:"documentId"`
	UserID     string           `json:"userId"`
	Passed     bool             `json:"passed"`
	Defects    []DocumentDefect `json:"defects"`
}
