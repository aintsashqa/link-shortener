package shortener

import (
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

type redirectService struct {
	redirectRepository RedirectRepositoryInterface
}

func NewRedirectService(redirectRepository RedirectRepositoryInterface) RedirectServiceInterface {
	return &redirectService{
		redirectRepository,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepository.Find(code)
}

func (r *redirectService) Add(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.Wrap(ErrRedirectInvalid, "service.Redirect.Add")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepository.Add(redirect)
}
