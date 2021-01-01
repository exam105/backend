package domain

type MetadataBson struct {
	Subject        string   `json:"subject,omitempty" bson:"subject,omitempty"`
	System         string   `json:"system,omitempty" bson:"system,omitempty"`
	Board          string   `json:"board,omitempty" bson:"board,omitempty"`
	Series         string   `json:"series,omitempty" bson:"series,omitempty"`
	Paper          string   `json:"paper,omitempty" bson:"paper,omitempty"`
	Year           string   `json:"year,omitempty" bson:"year,omitempty"`
	Month          string   `json:"month,omitempty" bson:"month,omitempty"`
	Username       string   `json:"username,omitempty" bson:"username"`
	Useremail      string   `json:"useremail,omitempty" bson:"useremail"`
	QuestionHexIds []string `json:"question_hex_ids,omitempty" bson:"question_hex_ids"`
}
