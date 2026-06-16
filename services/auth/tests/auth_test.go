package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	authv1 "github.com/antongolenev23/voltake-protos/gen/go/auth/v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	addr = "localhost:44044"
)

func newClient(t *testing.T) authv1.AuthClient {
	t.Helper()

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
	})

	return authv1.NewAuthClient(conn)
}

type User struct {
	Email    string
	Password string
}

func NewValidUser() User {
	return User{
		Email:    fmt.Sprintf("user_%s@test.com", gofakeit.LetterN(10)),
		Password: "StrongP@ssw0rd123!",
	}
}

func NewWeakUser() User {
	return User{
		Email:    fmt.Sprintf("user_%d@test.com", time.Now().UnixNano()),
		Password: "123",
	}
}

func NewInvalidEmailUser() User {
	return User{
		Email:    "not-an-email",
		Password: "StrongP@ssw0rd123!",
	}
}

func NewRandomEmail() string {
	return fmt.Sprintf("user_%s@test.com", gofakeit.LetterN(10))
}

func TestRegister_Success(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewValidUser()

	resp, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
}

func TestLogin_AfterRegister(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewValidUser()

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	resp, err := client.Login(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
}

func TestLogin_WithInvalidPassword(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewValidUser()

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	_, err = client.Login(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: fmt.Sprintf("%sabc", user.Password),
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, st.Code())
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewValidUser()

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	_, err = client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.AlreadyExists, st.Code())
}

func TestLogin_InvalidCredentials(t *testing.T) {
	t.Parallel()

	client := newClient(t)
	user := NewValidUser()

	_, err := client.Login(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, st.Code())
}

func TestRegister_InvalidEmail(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewInvalidEmailUser()

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_WeakPassword(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	user := NewWeakUser()

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    user.Email,
		Password: user.Password,
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_PasswordTooLong(t *testing.T) {
	t.Parallel()

	client := newClient(t)

	longPassword := gofakeit.Password(true, true, true, true, false, 80)

	require.Greater(t, len(longPassword), 72)

	_, err := client.Register(context.Background(), &authv1.Credentials{
		Email:    NewRandomEmail(),
		Password: longPassword,
	})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
