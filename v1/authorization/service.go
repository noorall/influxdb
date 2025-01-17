package authorization

import (
	"context"
	"time"

	"github.com/influxdata/influxdb/v2"
	"github.com/influxdata/influxdb/v2/kit/platform"
	"github.com/influxdata/influxdb/v2/kit/platform/errors"
	"github.com/influxdata/influxdb/v2/kv"
)

var (
	_ influxdb.AuthorizationService = (*Service)(nil)
	_ influxdb.PasswordsService     = (*Service)(nil)
)

type Service struct {
	store           *Store
	tenantService   TenantService
	strongPasswords bool
}

// NewService constructs a new Service.
func NewService(st *Store, ts TenantService, OptFns ...func(*Service)) *Service {
	svc := &Service{
		store:         st,
		tenantService: ts,
	}
	for _, fn := range OptFns {
		fn(svc)
	}
	return svc
}

func WithPasswordChecking(strong bool) func(*Service) {
	return func(s *Service) {
		s.strongPasswords = strong
	}
}

func (s *Service) CreateAuthorization(ctx context.Context, a *influxdb.Authorization) error {
	if err := a.Valid(); err != nil {
		return &errors.Error{
			Err: err,
		}
	}

	if a.Token == "" {
		return influxdb.ErrUnableToCreateToken
	}

	if _, err := s.tenantService.FindUserByID(ctx, a.UserID); err != nil {
		return influxdb.ErrUnableToCreateToken
	}

	if _, err := s.tenantService.FindOrganizationByID(ctx, a.OrgID); err != nil {
		return influxdb.ErrUnableToCreateToken
	}

	err := s.store.View(ctx, func(tx kv.Tx) error {
		if err := s.store.uniqueAuthToken(ctx, tx, a); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ErrTokenAlreadyExistsError
	}

	now := time.Now()
	a.SetCreatedAt(now)
	a.SetUpdatedAt(now)

	return s.store.Update(ctx, func(tx kv.Tx) error {
		return s.store.CreateAuthorization(ctx, tx, a)
	})
}

func (s *Service) FindAuthorizationByID(ctx context.Context, id platform.ID) (*influxdb.Authorization, error) {
	var a *influxdb.Authorization
	err := s.store.View(ctx, func(tx kv.Tx) error {
		auth, err := s.store.GetAuthorizationByID(ctx, tx, id)
		if err != nil {
			return err
		}

		a = auth
		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

// FindAuthorizationByToken returns a authorization by token for a particular authorization.
func (s *Service) FindAuthorizationByToken(ctx context.Context, n string) (*influxdb.Authorization, error) {
	var a *influxdb.Authorization
	err := s.store.View(ctx, func(tx kv.Tx) error {
		auth, err := s.store.GetAuthorizationByToken(ctx, tx, n)
		if err != nil {
			return err
		}

		a = auth

		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

// FindAuthorizations retrives all authorizations that match an arbitrary authorization filter.
// Filters using ID, or Token should be efficient.
// Other filters will do a linear scan across all authorizations searching for a match.
func (s *Service) FindAuthorizations(ctx context.Context, filter influxdb.AuthorizationFilter, opt ...influxdb.FindOptions) ([]*influxdb.Authorization, int, error) {
	if filter.ID != nil {
		var auth *influxdb.Authorization
		err := s.store.View(ctx, func(tx kv.Tx) error {
			a, e := s.store.GetAuthorizationByID(ctx, tx, *filter.ID)
			if e != nil {
				return e
			}
			auth = a
			return nil
		})
		if err != nil {
			return nil, 0, &errors.Error{
				Err: err,
			}
		}

		return []*influxdb.Authorization{auth}, 1, nil
	}

	if filter.Token != nil {
		var auth *influxdb.Authorization
		err := s.store.View(ctx, func(tx kv.Tx) error {
			a, e := s.store.GetAuthorizationByToken(ctx, tx, *filter.Token)
			if e != nil {
				return e
			}
			auth = a
			return nil
		})
		if err != nil {
			return nil, 0, &errors.Error{
				Err: err,
			}
		}

		return []*influxdb.Authorization{auth}, 1, nil
	}

	as := []*influxdb.Authorization{}
	err := s.store.View(ctx, func(tx kv.Tx) error {
		auths, err := s.store.ListAuthorizations(ctx, tx, filter)
		if err != nil {
			return err
		}
		as = auths
		return nil
	})

	if err != nil {
		return nil, 0, &errors.Error{
			Err: err,
		}
	}

	return as, len(as), nil
}

// UpdateAuthorization updates the status and description if available.
func (s *Service) UpdateAuthorization(ctx context.Context, id platform.ID, upd *influxdb.AuthorizationUpdate) (*influxdb.Authorization, error) {
	var auth *influxdb.Authorization
	err := s.store.View(ctx, func(tx kv.Tx) error {
		a, e := s.store.GetAuthorizationByID(ctx, tx, id)
		if e != nil {
			return e
		}
		auth = a
		return nil
	})

	if err != nil {
		return nil, &errors.Error{
			Code: errors.ENotFound,
			Err:  err,
		}
	}

	if upd.Status != nil {
		auth.Status = *upd.Status
	}
	if upd.Description != nil {
		auth.Description = *upd.Description
	}

	auth.SetUpdatedAt(time.Now())

	err = s.store.Update(ctx, func(tx kv.Tx) error {
		a, e := s.store.UpdateAuthorization(ctx, tx, id, auth)
		if e != nil {
			return e
		}
		auth = a
		return nil
	})
	return auth, err
}

func (s *Service) DeleteAuthorization(ctx context.Context, id platform.ID) error {
	return s.store.Update(ctx, func(tx kv.Tx) (err error) {
		return s.store.DeleteAuthorization(ctx, tx, id)
	})
}
