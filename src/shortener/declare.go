package shortener

import "errors"

var (
	ErrRedirectNotFound = errors.New("Redirect not found")
	ErrRedirectInvalid  = errors.New("Redirect invalid")
)

type Redirect struct {
	Code      string `json:"code" bson:"code"`
	Link      string `json:"link" bson:"link" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

type RedirectServiceInterface interface {
	Find(string) (*Redirect, error)
	Add(*Redirect) error
}

type RedirectRepositoryInterface interface {
	Find(string) (*Redirect, error)
	Add(*Redirect) error
}

type RedirectSerializerInterface interface {
	Decode([]byte) (*Redirect, error)
	Encode(*Redirect) ([]byte, error)
}
