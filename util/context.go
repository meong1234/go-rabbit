package util

import (
	"context"
)

const SessionKey = 1

type Session struct {
	CID    string
	Logger Logger
}

func SessionCid(ctx context.Context) string {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return ""
	}

	return session.CID
}

func SessionLogger(ctx context.Context) Logger {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return Log
	}

	return session.Logger
}

func NewSessionCtx(cid string, log Logger) context.Context {
	session := Session{
		cid,
		log,
	}
	return context.WithValue(context.Background(), SessionKey, &session)
}
