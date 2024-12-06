package service

import (
	"auth-service/cluster/userservice"
	"auth-service/domain"
	"context"
	"errors"
)

type Registration interface {
	SignUp(username, password string) (domain.TokenResponse, error)
}

type RegistrationServiceImpl struct {
	hasher      PasswordHasher
	jwtService  JwtService
	userService *userservice.Client
}

func NewRegistrationService(hasher PasswordHasher, jwtService JwtService, userClient *userservice.Client) Registration {
	return &RegistrationServiceImpl{
		hasher:      hasher,
		jwtService:  jwtService,
		userService: userClient,
	}
}

func (r *RegistrationServiceImpl) SignUp(username, password string) (domain.TokenResponse, error) {
	//сформировать хэш пароля
	passwordHash, err := r.hasher.Hash(password)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	//отправить пользователя на сохранение
	if userCheck, checkUserEerr := r.userService.GetUserByUserName(context.TODO(), username); checkUserEerr != nil || userCheck.Id != "" {
		return domain.TokenResponse{}, errors.New("user already exists")
	}

	user := domain.User{
		UserName: username,
		Password: passwordHash,
		Role:     domain.UserRole,
	}

	registeredUser, err := r.userService.RegisterUser(context.TODO(), user)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	//проверить ответ от сервиса пользователей
	if registeredUser.Id == "" {
		return domain.TokenResponse{}, errors.New("user registration failed")
	}

	//сгенерировать токен или отправить в кафку ивент о подтвержлениии пользователя
	//TODO: poka bez kafka
	accessToken, refreshToken, err := r.jwtService.Generate(registeredUser)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	//вернуть на фронт ответ с токенами или с текстом что пользователь ушел на подтверждение
	return domain.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
