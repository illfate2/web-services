package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	accessSecret []byte
}

func NewJWTService(accessSecret []byte) *JWTService {
	return &JWTService{
		accessSecret: accessSecret,
	}
}

func (s *JWTService) GenerateAccessToken(userID int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["uid"] = userID
	atClaims["access"] = true
	atClaims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString(s.accessSecret)
}

func (s *JWTService) GenerateRefreshToken(userID int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	atClaims["uid"] = userID
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString(s.accessSecret)
}

func (s *JWTService) RegenerateAccessToken(refreshToken string) (string, error) {
	jwtToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("no map claims")
	}

	userID, ok := claims["uid"].(int)
	if !ok {
		return "", errors.New("not valid user id type")
	}

	return s.GenerateAccessToken(userID)
}

func (s *JWTService) IsValidAccessToken(token string) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("no map claims")
	}

	access, ok := claims["access"].(bool)
	if !access || !ok {
		return errors.New("not valid token")
	}
	return nil
}

func (s *JWTService) UserIDFromToken(token string) (int, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("no map claims")
	}

	access, ok := claims["access"].(bool)
	if !access || !ok {
		return 0, errors.New("not valid token")
	}

	userID, ok := claims["uid"].(float64)
	if !ok {
		return 0, errors.New("not valid user id type")
	}

	exp := claims["exp"].(float64)
	isAfter := time.Unix(int64(exp), 0).After(time.Now())
	if !isAfter {
		return 0, errors.New("token expired")
	}
	return int(userID), nil
}
