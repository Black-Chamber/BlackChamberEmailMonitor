package m365

type M365Instance struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	APIMgmtURL   string
}

type MessageTrace struct {
	Organization     string  `json:"Organization"`
	MessageID        string  `json:"MessageId"`
	Received         string  `json:"Received"`
	SenderAddress    string  `json:"SenderAddress"`
	RecipientAddress string  `json:"RecipientAddress"`
	Subject          string  `json:"Subject"`
	Status           string  `json:"Status"`
	ToIP             *string `json:"ToIP"`   // Nullable field
	FromIP           *string `json:"FromIP"` // Nullable field
	Size             int     `json:"Size"`
	MessageTraceID   string  `json:"MessageTraceId"`
	StartDate        string  `json:"StartDate"`
	EndDate          string  `json:"EndDate"`
	Index            int     `json:"Index"`
}

type MessageTraceResponse struct {
	Value []MessageTrace `json:"value"` // Use the named MessageTrace type
}
