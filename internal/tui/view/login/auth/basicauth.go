package auth

import (
	"fmt"
	"strings"

	goauth "github.com/abbot/go-http-auth"
)

type BasicAuth struct {
	Users []string `cfg:"users"`
	Realm string   `cfg:"realm"`

	gAuth *goauth.BasicAuth `cfg:"-"`
}

func (b *BasicAuth) Prepare() error {
	v := make(map[string]string, len(b.Users))
	for _, user := range b.Users {
		parts := strings.Split(user, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid user: %q", user)
		}

		v[parts[0]] = parts[1]
	}

	if b.Realm == "" {
		b.Realm = "yap"
	}

	b.gAuth = &goauth.BasicAuth{
		Realm: b.Realm,
		Secrets: func(user, realm string) string {
			if hash, ok := v[user]; ok {
				return hash
			}

			return ""
		},
	}

	return nil
}

func (b *BasicAuth) Validator(username, password string) error {
	if b.gAuth == nil {
		return fmt.Errorf("basic auth not prepared")
	}

	if secret := b.gAuth.Secrets(username, b.Realm); secret == "" || !goauth.CheckSecret(password, secret) {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}
