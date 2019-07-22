package option

import "github.com/sextech/chaturbate/internal"

type ClientOption interface {
	Apply(*internal.Settings)
}

func WithoutAuthentication() ClientOption {
	return withoutAuthentication{}
}

type withoutAuthentication struct{}

func (w withoutAuthentication) Apply(o *internal.Settings) {
	o.NoAuth = true
}

func WithCredentials(username, password string) ClientOption {
	return withCredentials{
		username: username,
		password: password,
	}
}

type withCredentials struct {
	username string
	password string
}

func (w withCredentials) Apply(o *internal.Settings) {
	o.Username = w.username
	o.Password = w.password
}
