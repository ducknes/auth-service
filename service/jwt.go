package service

import (
	"auth-service/database"
	"auth-service/domain"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"strings"
	"time"
)

var (
	_accessTokenLifeTime = 15 * time.Minute
)

type JwtService interface {
	Generate(user domain.User) (string, string, error)
	Refresh(refreshToken string) (string, string, error)
	RetrieveRefreshToken(refreshToken string) error
}

type JwtServiceImpl struct {
	JwtSecret              string
	RefreshTokenRepository database.RefreshTokenRepository
}

func NewJwtService(secret string, repo database.RefreshTokenRepository) JwtService {
	return &JwtServiceImpl{
		JwtSecret:              secret,
		RefreshTokenRepository: repo,
	}
}

func (s *JwtServiceImpl) Generate(user domain.User) (string, string, error) {
	//сформировать клеймсы
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.UserName,
		"role":     user.Role,
		"exp":      time.Now().Add(_accessTokenLifeTime).Unix(),
	})

	//подписать аксес токен
	accessToken, err := claims.SignedString([]byte(s.JwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenBase := domain.RefreshTokenBase{
		UserId:   user.Id,
		UserName: user.UserName,
		UserRole: user.Role,
	}

	//сформировать рефереш токен
	refreshToken, err := s.generateRefreshToken(refreshTokenBase)
	if err != nil {
		return "", "", err
	}

	if err = s.RefreshTokenRepository.Set(context.TODO(), refreshToken); err != nil {
		return "", "", err
	}

	//вернуть 2 токена
	return accessToken, refreshToken, err
}

func (s *JwtServiceImpl) Refresh(refreshToken string) (string, string, error) {
	//получить из редиса рефреш токен
	redisToken, err := s.RefreshTokenRepository.Get(context.TODO(), refreshToken)
	if err != nil {
		return "", "", err
	}

	//проверить полученные данные
	userRefreshTokenBase, err := s.getRefreshTokenBase(refreshToken)
	if err != nil {
		return "", "", err
	}

	redisRefreshTokenBase, err := s.getRefreshTokenBase(redisToken)
	if err != nil {
		return "", "", err
	}

	if !userRefreshTokenBase.IsEquals(redisRefreshTokenBase) {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	//сгенерировать новые 2 токена
	return s.Generate(domain.User{Id: userRefreshTokenBase.UserId, UserName: userRefreshTokenBase.UserName, Role: userRefreshTokenBase.UserRole})
}

func (s *JwtServiceImpl) RetrieveRefreshToken(refreshToken string) error {
	return s.RefreshTokenRepository.Delete(context.TODO(), refreshToken)
}

func (s *JwtServiceImpl) generateRefreshToken(tokenBase domain.RefreshTokenBase) (string, error) {
	baseBytes, err := json.Marshal(tokenBase)
	if err != nil {
		return "", err
	}

	refreshTokenSalt, err := randomBytesString()
	if err != nil {
		return "", err
	}

	base := base64.URLEncoding.EncodeToString(baseBytes)

	return fmt.Sprintf("%s.%s", base, refreshTokenSalt), nil
}

func randomBytesString() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *JwtServiceImpl) getRefreshTokenBase(token string) (base domain.RefreshTokenBase, err error) {
	if len(token) == 0 {
		return domain.RefreshTokenBase{}, errors.New("token is empty")
	}

	b, err := base64.StdEncoding.DecodeString(strings.Split(token, ".")[0])
	if err != nil {
		return domain.RefreshTokenBase{}, err
	}

	return base, json.Unmarshal(b, &base)
}
