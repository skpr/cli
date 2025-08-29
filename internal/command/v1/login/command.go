package login

import (
	"context"
	"fmt"
	"log"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/aws/cognito/oidc/login"
	"github.com/skpr/cli/internal/aws/cognito/oidc/rand"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config/clusters"
	credentialscache "github.com/skpr/cli/internal/client/config/credentials/cache"
	skprdiscovery "github.com/skpr/cli/internal/client/config/discovery"
	"github.com/skpr/cli/internal/client/config/project"
)

const (
	// StateLength is the length of the state string used to prevent CSRF.
	StateLength = 8
)

// Command to login to the platform.
type Command struct {
	Callback string
}

// Run the login command.
func (cmd *Command) Run() error {
	log.Println("Connecting to cluster")

	discovery, err := skprdiscovery.New()
	if err != nil {
		return fmt.Errorf("failed to discover config %w", err)
	}

	configFile, err := discovery.Config()
	if err != nil {
		return fmt.Errorf("failed to get project config file path %w", err)
	}

	config, err := project.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to get project config %w", err)
	}

	clusterFile, err := discovery.Clusters()
	if err != nil {
		return fmt.Errorf("failed to get clusters file %w", err)
	}

	clusterCfg, err := clusters.LoadFromFile(clusterFile, config.Cluster)
	if err != nil {
		return fmt.Errorf("failed to get clusters config %w", err)
	}

	conn, err := client.Dial(clusterCfg.API)
	if err != nil {
		return err
	}

	// @todo, We should not have to do this many function calls
	//        to setup a client.

	client := pb.NewLoginClient(conn)

	ctx := context.TODO()

	providerInfo, err := client.GetProviderInfo(ctx, &pb.LoginGetProviderInfoRequest{})
	if err != nil {
		return err
	}

	if providerInfo.Cognito == nil {
		return fmt.Errorf("unknown login provider")
	}

	log.Println("Found login provider information")

	// State is used to prevent CSRF, this is a random string.
	// It is attached to the request and returned in the response.
	state := rand.String(StateLength)

	ctxReady, ready := context.WithCancel(context.Background())

	server := login.NewServer(cmd.Callback)

	group, _ := errgroup.WithContext(context.Background())

	oath2Config := oauth2.Config{
		RedirectURL: cmd.Callback,
		ClientID:    providerInfo.Cognito.ClientID,
		// @todo, These should be managed at our API layer.
		// https://previousnext.atlassian.net/browse/SKPR-1001
		Scopes: []string{"openid email profile aws.cognito.signin.user.admin"},
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   providerInfo.Cognito.AuthURL,
			TokenURL:  providerInfo.Cognito.TokenURL,
		},
	}

	group.Go(func() error {
		log.Println("Starting webserver for login callback")

		resp, err := server.Run(context.TODO(), ready)
		if err != nil {
			fmt.Println("Failed to start server:", err)
		}

		log.Println("Callback received")

		if resp.Error != "" {
			return fmt.Errorf("login failed: error_code=%s error_description=%s", resp.Code, resp.ErrorDescription)
		}

		// Ensure that we are secure.
		if resp.State != state {
			return fmt.Errorf("failed to login with code: %w", err)
		}

		token, err := oath2Config.Exchange(context.TODO(), resp.Code)
		if err != nil {
			return fmt.Errorf("failed to login with code: %w", err)
		}

		// Store this information for later so we can generate temporary credentials.
		credentials := credentialscache.Credentials{
			Config: oath2Config,
			Token: credentialscache.Token{
				Refresh: token.RefreshToken,
			},
			Cognito: credentialscache.Cognito{
				Region:             providerInfo.Cognito.Region,
				IdentityPoolID:     providerInfo.Cognito.IdentityPoolID,
				IdentityProviderID: providerInfo.Cognito.IdentityProviderID,
			},
		}

		err = credentialscache.Set(config.Cluster, credentials)
		if err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctxReady.Done()

		log.Println("Opening browser session")

		authURL := oath2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

		fmt.Println("Authentication URL:", authURL)

		return open.Run(authURL)
	})

	err = group.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for login: %w", err)
	}

	log.Println("Successfully stored temporary credentials")

	return nil
}
