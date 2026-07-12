package authclient

import (
	"context"

	authv1 "github.com/antongolenev23/voltake-protos/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/antongolenev23/voltake-services/services/booking/internal/auth-client/interceptor"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	JWTToken string `json:"jwtToken"`
}

type Client struct {
	conn   *grpc.ClientConn
	client authv1.AuthClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			interceptor.RequestID(),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: authv1.NewAuthClient(conn),
	}, nil
}

func (c *Client) Register(
	ctx context.Context,
	credentials Credentials,
) (AuthResponse, error) {

	resp, err := c.client.Register(ctx, &authv1.Credentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	})

	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		JWTToken: resp.GetToken(),
	}, nil
}

func (c *Client) Login(
	ctx context.Context,
	credentials Credentials,
) (AuthResponse, error) {

	resp, err := c.client.Login(ctx, &authv1.Credentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	})

	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		JWTToken: resp.GetToken(),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
