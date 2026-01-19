package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type UserData struct {
	UserID uint
	Role   string
}

func GetUserData(ctx context.Context) (UserData, error) {
	if v := ctx.Value("claims"); v != nil {
		uintUserID, err := strconv.ParseUint(v.(Claims).UserID, 10, 64)
		if err != nil {
			return UserData{}, err
		}

		return UserData{
			UserID: uint(uintUserID),
			Role:   v.(Claims).Role,
		}, nil
	}

	return UserData{}, fmt.Errorf("unable to get data from context")
}
