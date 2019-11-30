package notify

import (
	"context"
	"errors"
	"net/http"
)

type key int

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
const (
	notifyContextKey     key = iota
)

// NotifyFromRequest retrieves the current Notifier associated with the request, error is set appropriately
func NotifyFromRequest(w http.ResponseWriter, r *http.Request) (*Notifier, error) {
	ctx := r.Context()
	return NotifyFromContext(ctx)
}

// NotifyFromContext gets the current Notifier from the given context
func NotifyFromContext(ctx context.Context) (*Notifier, error) {

	notify, ok := ctx.Value(notifyContextKey).(*Notifier)
	if ok && notify != nil {
		return notify, nil
	}

	return nil, errors.New("no notifier in context")
}

