package processor

type DetectionRules struct {
	Service string `yaml:"service"`
	Rules   []Rule `yaml:"rules"`
}

type Rule struct {
	Name       string      `yaml:"name"`
	Conditions []Condition `yaml:"conditions"`
	Confidence int         `yaml:"confidence"`
}

type Condition struct {
	Field string `yaml:"field"`
	Regex string `yaml:"regex"`
}

type Message struct {
	Organization     string
	MessageID        string
	Received         string
	SenderAddress    string
	RecipientAddress string
	Subject          string
	Status           string
	ToIP             *string
	FromIP           *string
	Size             int
	MessageTraceID   string
	StartDate        string
	EndDate          string
	Index            int
}
