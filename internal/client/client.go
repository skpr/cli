package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skpr/api/pb"
	"github.com/skpr/cli/internal/client/config/clusters"
	skprdiscovery "github.com/skpr/cli/internal/client/config/discovery"
	"github.com/skpr/cli/internal/client/config/project"
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

// Options for the client.
type Options struct {
	Username      string
	Password      string
	Credentials   string
	ClusterConfig string
}

// New client built using options.
func (o *Options) New() (*Client, context.Context, error) {
	return NewFromFile()
}

// Client for interacting with the Skipper server.
type Client struct {
	ClientConn          *grpc.ClientConn
	config              project.Config
	Discovery           *skprdiscovery.Discovery
	CredentialsProvider aws.CredentialsProvider
	cluster             clusters.Cluster
}

// New client.
func New() (*Client, context.Context, error) {
	dial, err := Dial()
	return &Client{ClientConn: dial}, context.TODO(), err
}

// Dial a connection to the API server.
func Dial() (*grpc.ClientConn, error) {
	return grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewFromFile loads a file and uses that configuration to return a client.
func NewFromFile() (*Client, context.Context, error) {
	return New()
}

// Project client.
func (c Client) Project() pb.ProjectClient {
	return pb.NewProjectClient(c.ClientConn)
}

// Environment client operations.
func (c Client) Environment() pb.EnvironmentClient {
	return pb.NewEnvironmentClient(c.ClientConn)
}

// Config client operations.
func (c Client) Config() pb.ConfigClient {
	return pb.NewConfigClient(c.ClientConn)
}

// Cron client operations.
func (c Client) Cron() pb.CronClient {
	return pb.NewCronClient(c.ClientConn)
}

// Purge client operations.
func (c Client) Purge() pb.PurgeClient {
	return pb.NewPurgeClient(c.ClientConn)
}

// Login client operations.
func (c Client) Login() pb.LoginClient {
	return pb.NewLoginClient(c.ClientConn)
}

// Logs client operations.
func (c Client) Logs() pb.LogsClient {
	return pb.NewLogsClient(c.ClientConn)
}

// Backup client operations.
func (c Client) Backup() pb.BackupClient {
	return pb.NewBackupClient(c.ClientConn)
}

// Release client operations.
func (c Client) Release() pb.ReleaseClient {
	return pb.NewReleaseClient(c.ClientConn)
}

// Restore client operations.
func (c Client) Restore() pb.RestoreClient {
	return pb.NewRestoreClient(c.ClientConn)
}

// Image client operations.
func (c Client) Image() pb.MysqlClient {
	return pb.NewMysqlClient(c.ClientConn)
}

// Mysql client operations.
func (c Client) Mysql() pb.MysqlClient {
	return pb.NewMysqlClient(c.ClientConn)
}

// SSH client operations.
func (c Client) SSH() ssh.Interface {
	return ssh.Client{Config: c.config, CredentialsProvider: c.CredentialsProvider, Cluster: c.cluster}
}

// Version client operations.
func (c Client) Version() pb.VersionClient {
	return pb.NewVersionClient(c.ClientConn)
}

// Volume client operations.
func (c Client) Volume() pb.VolumeClient {
	return pb.NewVolumeClient(c.ClientConn)
}

// Daemon client operations.
func (c Client) Daemon() pb.DaemonClient {
	return pb.NewDaemonClient(c.ClientConn)
}

// Compass client operations.
func (c Client) Compass() pb.CompassClient {
	return pb.NewCompassClient(c.ClientConn)
}
