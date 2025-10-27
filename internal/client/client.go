package client

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/skpr/cli/internal/client/config"
	skprcredentials "github.com/skpr/cli/internal/client/credentials"
	"github.com/skpr/cli/internal/client/ssh"
)

const (
	// KeyProject is used for determining a project name from the metadata.
	KeyProject = "project"
	// KeyUsername is used for determining user credentials on the platform.
	KeyUsername = "username"
	// KeyPassword is used for determining user credentials on the platform.
	KeyPassword = "password"
	// KeySession is used for determining user credentials on the platform.
	KeySession = "session"
)

// Client for interacting with the Skipper server.
type Client struct {
	conn *grpc.ClientConn
	// These are used by the ssh client.
	config config.Config
	// This is public so other clients eg. package and utilise them.
	Credentials skprcredentials.Credentials
}

// New client.
func New(ctx context.Context) (context.Context, *Client, error) {
	config, err := config.New()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create config: %w", err)
	}

	credentials, err := skprcredentials.New(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not retrieve credentials: %w", err)
	}

	conn, err := Dial(config)
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to dial server: %w", err)
	}

	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
	md := metadata.Pairs(
		KeyProject, config.Project,
		KeyUsername, credentials.Username,
		KeyPassword, credentials.Password,
		KeySession, credentials.Session,
	)

	client := &Client{
		conn: conn,
		// These are used by the ssh client.
		config: config,
		// This is public so other clients eg. package and utilise them.
		Credentials: credentials,
	}

	return metadata.NewOutgoingContext(ctx, md), client, err
}

// Dial a connection to the API server.
func Dial(config config.Config) (*grpc.ClientConn, error) {
	server := fmt.Sprintf("%s:%d", config.API.Host(), config.API.Port())

	if config.API.Insecure() {
		return grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return grpc.NewClient(server, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
}

// Project client.
func (c Client) Project() pb.ProjectClient {
	return pb.NewProjectClient(c.conn)
}

// Environment client operations.
func (c Client) Environment() pb.EnvironmentClient {
	return pb.NewEnvironmentClient(c.conn)
}

// Config client operations.
func (c Client) Config() pb.ConfigClient {
	return pb.NewConfigClient(c.conn)
}

// Cron client operations.
func (c Client) Cron() pb.CronClient {
	return pb.NewCronClient(c.conn)
}

// Purge client operations.
func (c Client) Purge() pb.PurgeClient {
	return pb.NewPurgeClient(c.conn)
}

// Login client operations.
func (c Client) Login() pb.LoginClient {
	return pb.NewLoginClient(c.conn)
}

// Logs client operations.
func (c Client) Logs() pb.LogsClient {
	return pb.NewLogsClient(c.conn)
}

// Backup client operations.
func (c Client) Backup() pb.BackupClient {
	return pb.NewBackupClient(c.conn)
}

// Release client operations.
func (c Client) Release() pb.ReleaseClient {
	return pb.NewReleaseClient(c.conn)
}

// Restore client operations.
func (c Client) Restore() pb.RestoreClient {
	return pb.NewRestoreClient(c.conn)
}

// Image client operations.
func (c Client) Image() pb.MysqlClient {
	return pb.NewMysqlClient(c.conn)
}

// Mysql client operations.
func (c Client) Mysql() pb.MysqlClient {
	return pb.NewMysqlClient(c.conn)
}

// SSH client operations.
func (c Client) SSH() ssh.Interface {
	return ssh.Client{Config: c.config, Credentials: c.Credentials}
}

// Version client operations.
func (c Client) Version() pb.VersionClient {
	return pb.NewVersionClient(c.conn)
}

// Volume client operations.
func (c Client) Volume() pb.VolumeClient {
	return pb.NewVolumeClient(c.conn)
}

// Daemon client operations.
func (c Client) Daemon() pb.DaemonClient {
	return pb.NewDaemonClient(c.conn)
}
