package main

import (
	"net/http"
	"net/url"
	"time"
	"fmt"
)

const (
	sessionLength     = 24 * 3 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLength   = 20
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)
	session := &Session{
		ID:     GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}
	http.SetCookie(w, &cookie)
	return session
}

func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}
	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		return nil
	}
	if session == nil {
		return nil
	}
	if session.IsExpired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserID == "" {
		return nil
	}
	user, err := globalUserStore.Find(session.UserID)
	if err != nil {
		return nil // just redirect or something.
	}
	return user
}

func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if RequestUser(r) != nil {
		fmt.Println("has an existing session")
		return
	}
	fmt.Println("need to login!")
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}
	return session
}
