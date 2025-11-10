package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gabrielfmcoelho/platform-core/internal/password"
)

type UserUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *UserUsecase) Create(c context.Context, createUser *domain.CreateUser) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	_, err := uu.userRepository.GetByEmail(ctx, createUser.Email)
	if err == nil {
		return domain.ErrUserAlreadyExists
	}

	hashedPassword, err := password.HashPassword(createUser.Password)
	if err != nil {
		return err
	}
	createUser.Password = hashedPassword

	user := parser.ToUser(createUser)

	err = uu.userRepository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (uu *UserUsecase) Fetch(c context.Context) ([]domain.PublicUser, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	users, err := uu.userRepository.Fetch(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	// parse the users to domain.PublicUser
	publicUsers := make([]domain.PublicUser, 0)
	for _, user := range users {
		publicUsers = append(publicUsers, parser.ToPublicUser(user))
	}

	return publicUsers, nil
}

func (uu *UserUsecase) GetByIdentifier(c context.Context, identifier string) (domain.PublicUser, error) {
	// if the identifier is an email, get the user by email
	// if the identifier is an ID, get the user by ID
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	var user domain.User
	var publicUser domain.PublicUser
	var err error
	var id uint
	if len(identifier) > 0 {
		if identifier[0] >= '0' && identifier[0] <= '9' {
			id, err = internal.ParseUint(identifier)
			if err != nil {
				return publicUser, err
			}
			user, err = uu.userRepository.GetByID(ctx, id)
		} else {
			user, err = uu.userRepository.GetByEmail(ctx, identifier)
		}
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return publicUser, domain.ErrNotFound
			}
			return publicUser, domain.ErrInternalServerError
		}
	} else {
		return publicUser, domain.ErrInvalidIdentifier
	}
	return parser.ToPublicUser(user), nil
}

func (uu *UserUsecase) Update(c context.Context, userID uint, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	err := uu.userRepository.Update(ctx, userID, user)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (uu *UserUsecase) Archive(c context.Context, userID uint) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	err := uu.userRepository.Archive(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (uu *UserUsecase) Unarchive(c context.Context, userID uint) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	err := uu.userRepository.Unarchive(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}
