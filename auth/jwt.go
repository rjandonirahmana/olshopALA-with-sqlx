package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(customer_id int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GenerateTokenSeller(seller_id int) (string, error)
	ValidateTokenSeller(encodedToken string) (*jwt.Token, error)
}

type jwtService struct {
	SECRET_KEY  []byte
	SECRET_KEY2 []byte
}

func NewService(secret_key_customer, secret_key_seller string) *jwtService {
	return &jwtService{SECRET_KEY: []byte(secret_key_customer), SECRET_KEY2: []byte(secret_key_seller)}
}

func (j *jwtService) GenerateToken(customer_id int) (string, error) {

	//claim adalah payload data jwt
	claim := jwt.MapClaims{}
	claim["customer_id"] = customer_id
	claim["exp"] = time.Now().Add(time.Hour * 8).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(j.SECRET_KEY)

	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (j *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("INVALID ERROR")
		}

		return []byte(j.SECRET_KEY), nil

	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func (j *jwtService) GenerateTokenSeller(seller_id int) (string, error) {

	//claim adalah payload data jwt
	claim := jwt.MapClaims{}
	claim["seller_id"] = seller_id
	claim["exp"] = time.Now().Add(time.Hour * 8).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(j.SECRET_KEY2)

	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (j *jwtService) ValidateTokenSeller(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("INVALID ERROR")
		}

		return []byte(j.SECRET_KEY2), nil

	})
	if err != nil {
		return token, err
	}

	return token, nil

}
