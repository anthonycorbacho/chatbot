package bot

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Bot represents a very limited bot that does not understand much,
// but tried his best to answer question he already knows the answer.
type Bot struct {

	// brian keep and the mapping between precomputed regex and answer function.
	brain map[*regexp.Regexp]HandlerFunc

	// We can extend this a little bit and can in the future add a storage (in mem or redis)
	// so user can ask bit to remember something and ask later what was it.
	// Example:
	//	bot.New(redis.Memory("localhost:6379"))
	//	bot.Respond("remember (.+) is (.+)", b.Remember)
	//	bot.Respond("what is (.+)", b.WhatIs)
	//
	// and the function can be something like this
	//
	//	func (b *Bot) Remember(msg Message) string {
	//		key, value := msg.Matches[0], msg.Matches[1]
	//		b.Store.Set(key, value)
	//		return "OK, I'll remember %s is %s", key, value
	//	}
	//
	//	func (b *Bot) WhatIs(msg Message) string {
	//		key := msg.Matches[0]
	//		var value string
	//		ok, err := b.Store.Get(key, &value)
	//		if err != nil {
	//			return "Humm, cannot remember, i lost my brain."
	//		}
	//		if !ok {
	//			return "I do not remember %q", key
	//		}
	//		return "%s is %s", key, value
	//	}
}

// HandlerFunc represents the function signature required for implementing a bot response.
type HandlerFunc = func(Message) string

// Message is automatically created via Respond or RespondRegex
// when the message matches the regular expression of the handler.
type Message struct {
	Context context.Context
	Text    string
	Matches []string
}

// New creates a new Bot and initializes it.
func New() *Bot {
	return &Bot{
		brain: make(map[*regexp.Regexp]HandlerFunc),
	}
}

// Respond registers and executes the given function only if the message
// text matches the given message.
// The message will be matched against the msg string as regular
// expression that must match the entire message in a case insensitive way.
//
// You can use sub matches in the msg which will be passed to the function via
// Matches.
//
// If you need complete control over the regular expression, e.g. because you
// want the patter to match only a substring of the message but not all of it,
// you can use RespondRegex.
//
func (b *Bot) Respond(pattern string, f HandlerFunc) error {
	expr := fmt.Sprintf("^%s$", pattern)
	return b.RespondRegex(expr, f)
}

// RespondRegex is like Respond but gives a little more control over the
// regular expression. However, also with this function messages are matched in
// a case insensitive way.
func (b *Bot) RespondRegex(expr string, f HandlerFunc) error {
	if expr == "" {
		return ErrEmptyPattern
	}

	if expr[0] == '^' {
		// String starts with the "^" anchor but does it also have the prefix
		// or case insensitive matching?
		if !strings.HasPrefix(expr, "^(?i)") { // TODO: strings.ToLower would be easier?
			expr = "^(?i)" + expr[1:]
		}
	} else {
		// The string is not starting with "^" but maybe it has the prefix for
		// case insensitive matching already?
		if !strings.HasPrefix(expr, "(?i)") {
			expr = "(?i)" + expr
		}
	}

	regex, err := regexp.Compile(expr)
	if err != nil {
		return fmt.Errorf("could not compile pattern %s: %w", expr, err)
	}
	b.brain[regex] = f
	return nil
}

// Sentence takes a input text or a question and will match again boot knowledge.
// if no match found a generic sorry message will be returned.
func (b *Bot) Sentence(ctx context.Context, msg string) string {

	// This is clairely not the best option, the runtime complexity is O(n)
	// this can be improved, I did not have time to reseach how chatbot should be
	// implemented, so I went to the simple solution.
	for k, v := range b.brain {
		matches := k.FindStringSubmatch(msg)
		if len(matches) == 0 {
			continue
		}
		return v(Message{
			Context: ctx,
			Matches: matches[1:],
			Text:    msg,
		})
	}
	return "Sorry, I do not know the answer, I will guess the answer will be 42."
}
