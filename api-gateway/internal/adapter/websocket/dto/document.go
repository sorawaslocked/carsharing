package wsdto

type DocumentDefect struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type DocumentAnalyzedMessage struct {
	DocumentID string           `json:"documentID"`
	UserID     string           `json:"userID"`
	Passed     bool             `json:"passed"`
	Defects    []DocumentDefect `json:"defects"`
}
