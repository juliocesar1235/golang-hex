package shortener

import (
	"fmt"
	"errors"
	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

func NewRedirectServiceInstance(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (rService *redirectService) Find(code string) (*Redirect, error) {
	return rService.redirectRepo.Find(code)
}

func (rService *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}


	fmt.Println("RedirectSERVICE ", redirect)
	redirect.Code = shortid.MustGenerate()
	fmt.Println("RedirectCODE ", redirect)
	redirect.CreatedAt = time.Now().UTC().Unix()
	fmt.Println("RedirectCREATEDAT ", redirect)
	return rService.redirectRepo.Store(redirect)
}