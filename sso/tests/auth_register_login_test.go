package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssov1 "github.com/tizzhh/auth-grpc-service/protos/gen/go/sso"
	"github.com/tizzhh/auth-grpc-service/sso/tests/suite"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passwordDefaultLen = 10
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with empty email and password",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Loign with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with empty email and password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with non-matching password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid email and/or password",
		},
		{
			name:        "Login without appID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passwordDefaultLen)
}
