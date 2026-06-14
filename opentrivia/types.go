package opentrivia

// Question is one trivia question from the Open Trivia Database.
type Question struct {
	Rank          int    `json:"rank"`
	Category      string `kit:"id" json:"category"`
	Type          string `json:"type"`           // "multiple" or "boolean"
	Difficulty    string `json:"difficulty"`     // "easy", "medium", "hard"
	Question      string `json:"question"`       // HTML-decoded
	CorrectAnswer string `json:"correct_answer"` // HTML-decoded
	OtherAnswers  string `json:"other_answers"`  // comma-joined incorrect answers, HTML-decoded
}

// Category is one entry from the categories endpoint.
type Category struct {
	Rank int    `json:"rank"`
	ID   int    `kit:"id" json:"id"`
	Name string `json:"name"`
}

// Token is the result from the session token endpoint.
type Token struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// unexported: used inside opentrivia.go for JSON decode only.

type apiResponse struct {
	ResponseCode int           `json:"response_code"`
	Results      []apiQuestion `json:"results"`
}

type apiQuestion struct {
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type categoryResponse struct {
	TriviaCategories []apiCategory `json:"trivia_categories"`
}

type apiCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type tokenResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Token           string `json:"token"`
}
