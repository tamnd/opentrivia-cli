// Package opentrivia exposes the Open Trivia Database as a kit Domain driver.
// A multi-domain host (ant) enables it with a single blank import:
//
//	import _ "github.com/tamnd/opentrivia-cli/opentrivia"
//
// The same Domain also builds the standalone opentrivia binary (see cli/root.go),
// so the binary and a host share one source of truth.
package opentrivia

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the opentrivia driver.
type Domain struct{}

// Info describes the scheme, the hostnames a pasted link is matched against,
// and the identity reused for the binary's help and version.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "opentrivia",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "opentrivia",
			Short:  "Open Trivia Database — free trivia questions",
			Long: `opentrivia fetches trivia questions, category lists, and session tokens
from opentdb.com. No API key required.`,
			Site: Host,
			Repo: "https://github.com/tamnd/opentrivia-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	kit.Handle(app, kit.OpMeta{
		Name:    "questions",
		Group:   "read",
		List:    true,
		Summary: "Fetch trivia questions",
	}, questionsOp)

	kit.Handle(app, kit.OpMeta{
		Name:    "categories",
		Group:   "read",
		List:    true,
		Summary: "List all trivia categories",
	}, categoriesOp)

	kit.Handle(app, kit.OpMeta{
		Name:    "token",
		Group:   "read",
		Single:  true,
		Summary: "Request a new session token",
	}, tokenOp)
}

// newClient builds the client from host-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type questionsInput struct {
	Amount     int     `kit:"flag,inherit" help:"number of questions (1-50)"`
	Category   int     `kit:"flag,inherit" help:"category ID (0 = any)"`
	Difficulty string  `kit:"flag,inherit" help:"difficulty: easy, medium, hard"`
	Type       string  `kit:"flag,inherit" help:"question type: multiple, boolean"`
	Token      string  `kit:"flag,inherit" help:"session token to avoid repeat questions"`
	Client     *Client `kit:"inject"`
}

type categoriesInput struct {
	Client *Client `kit:"inject"`
}

type tokenInput struct {
	Client *Client `kit:"inject"`
}

// --- handlers ---

func questionsOp(ctx context.Context, in questionsInput, emit func(Question) error) error {
	amount := in.Amount
	if amount <= 0 {
		amount = 10
	}
	questions, err := in.Client.Questions(ctx, amount, in.Category, in.Difficulty, in.Type, in.Token)
	if err != nil {
		return mapErr(err)
	}
	for _, q := range questions {
		if err := emit(q); err != nil {
			return err
		}
	}
	return nil
}

func categoriesOp(ctx context.Context, in categoriesInput, emit func(Category) error) error {
	cats, err := in.Client.Categories(ctx)
	if err != nil {
		return mapErr(err)
	}
	for _, cat := range cats {
		if err := emit(cat); err != nil {
			return err
		}
	}
	return nil
}

func tokenOp(ctx context.Context, in tokenInput, emit func(Token) error) error {
	tok, err := in.Client.Token(ctx)
	if err != nil {
		return mapErr(err)
	}
	return emit(*tok)
}

// --- Resolver ---

// Classify turns an input into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty opentrivia reference")
	}
	return "questions", input, nil
}

// Locate returns the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "questions":
		return "https://opentdb.com/", nil
	case "categories":
		return "https://opentdb.com/api_category.php", nil
	default:
		return "", errs.Usage("opentrivia has no resource type %q", uriType)
	}
}

// mapErr converts a library error into the kit error kind.
func mapErr(err error) error {
	return err
}
