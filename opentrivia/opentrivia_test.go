package opentrivia_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/opentrivia-cli/opentrivia"
)

const fakeQuestionsJSON = `{
  "response_code": 0,
  "results": [
    {
      "type": "multiple",
      "difficulty": "medium",
      "category": "Science: Computers",
      "question": "What does &quot;HTTP&quot; stand for?",
      "correct_answer": "HyperText Transfer Protocol",
      "incorrect_answers": [
        "High Transfer Text Protocol",
        "HyperText Transmission Protocol",
        "HyperType Transfer Protocol"
      ]
    },
    {
      "type": "boolean",
      "difficulty": "easy",
      "category": "Science: Computers",
      "question": "Linux was first created as an alternative to Windows XP.",
      "correct_answer": "False",
      "incorrect_answers": ["True"]
    }
  ]
}`

const fakeCategoriesJSON = `{
  "trivia_categories": [
    {"id": 9, "name": "General Knowledge"},
    {"id": 10, "name": "Entertainment: Books"},
    {"id": 11, "name": "Entertainment: Film"}
  ]
}`

const fakeTokenJSON = `{
  "response_code": 0,
  "response_message": "Token Generated Successfully!",
  "token": "f00cb469ce38726ee00a7c6836761b0a4fb808181a125dcde6d50a9f3c9127b6"
}`

const fakeNoResultsJSON = `{
  "response_code": 1,
  "results": []
}`

func newTestClient(ts *httptest.Server) *opentrivia.Client {
	cfg := opentrivia.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return opentrivia.NewClient(cfg)
}

func TestQuestionsSendsUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = fmt.Fprint(w, fakeQuestionsJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Questions(context.Background(), 2, 0, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if gotUA == "" {
		t.Error("User-Agent not sent")
	}
}

func TestQuestionsParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeQuestionsJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.Questions(context.Background(), 2, 0, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].Rank != 1 {
		t.Errorf("items[0].Rank = %d, want 1", items[0].Rank)
	}
	if items[0].Category != "Science: Computers" {
		t.Errorf("items[0].Category = %q", items[0].Category)
	}
	if items[0].Difficulty != "medium" {
		t.Errorf("items[0].Difficulty = %q", items[0].Difficulty)
	}
	if items[0].CorrectAnswer != "HyperText Transfer Protocol" {
		t.Errorf("items[0].CorrectAnswer = %q", items[0].CorrectAnswer)
	}
	wantOther := "High Transfer Text Protocol, HyperText Transmission Protocol, HyperType Transfer Protocol"
	if items[0].OtherAnswers != wantOther {
		t.Errorf("items[0].OtherAnswers = %q, want %q", items[0].OtherAnswers, wantOther)
	}
}

func TestQuestionsHTMLDecoded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeQuestionsJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.Questions(context.Background(), 2, 0, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	want := `What does "HTTP" stand for?`
	if items[0].Question != want {
		t.Errorf("items[0].Question = %q, want %q", items[0].Question, want)
	}
}

func TestQuestionsRetriesOn503(t *testing.T) {
	var hits int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprint(w, fakeQuestionsJSON)
	}))
	defer ts.Close()

	cfg := opentrivia.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 3
	c := opentrivia.NewClient(cfg)

	_, err := c.Questions(context.Background(), 2, 0, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
}

func TestQuestionsNoResults(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeNoResultsJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Questions(context.Background(), 5, 0, "", "", "")
	if err == nil {
		t.Error("expected error for response_code=1, got nil")
	}
}

func TestCategoriesParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeCategoriesJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	cats, err := c.Categories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 3 {
		t.Fatalf("len(cats) = %d, want 3", len(cats))
	}
	if cats[0].Rank != 1 {
		t.Errorf("cats[0].Rank = %d, want 1", cats[0].Rank)
	}
	if cats[0].ID != 9 {
		t.Errorf("cats[0].ID = %d, want 9", cats[0].ID)
	}
	if cats[0].Name != "General Knowledge" {
		t.Errorf("cats[0].Name = %q, want %q", cats[0].Name, "General Knowledge")
	}
}

func TestTokenParses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeTokenJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tok, err := c.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token == "" {
		t.Error("Token is empty")
	}
	if tok.Message == "" {
		t.Error("Message is empty")
	}
	want := "f00cb469ce38726ee00a7c6836761b0a4fb808181a125dcde6d50a9f3c9127b6"
	if tok.Token != want {
		t.Errorf("Token = %q, want %q", tok.Token, want)
	}
}
