package logout

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/api/pb"
	oidclogin "github.com/skpr/cli/internal/aws/cognito/oidc/login"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config"
	credentialscache "github.com/skpr/cli/internal/client/credentials/cache"
)

// Command to logout from the platform.
type Command struct {
	Callback string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	log.Println("Connecting to cluster")

	config, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to get project config: %w", err)
	}

	credentials, found, err := credentialscache.Get(config.Cluster)
	if err != nil {
		return fmt.Errorf("failed to get cached credentials: %w", err)
	}

	if found {
		log.Println("Deleting cached credentials")

		token, err := credentials.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get token source: %w", err)
		}

		cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(credentials.Cognito.Region), awsconfig.WithCredentialsProvider(aws.AnonymousCredentials{}))
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		_, err = cognitoidentityprovider.NewFromConfig(cfg).GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
			AccessToken: aws.String(token.AccessToken),
		})
		if err != nil {
			// @todo, We can add this later once we have log levels eg. Turn of debug mode.
			// https://previousnext.atlassian.net/browse/SKPR-1002
			fmt.Println("Failed to execute global sign out:", err)
		}

		// Delete the file now that we have invalidated our tokens.
		err = credentialscache.Delete(config.Cluster)
		if err != nil {
			return fmt.Errorf("failed to delete credentials cache %w", err)
		}
	}

	conn, err := client.Dial(config)
	if err != nil {
		return fmt.Errorf("could not connect to cluster: %w", err)
	}

	loginClient := pb.NewLoginClient(conn)

	providerInfo, err := loginClient.GetProviderInfo(ctx, &pb.LoginGetProviderInfoRequest{})
	if err != nil {
		return err
	}

	if providerInfo.Cognito == nil {
		return fmt.Errorf("unknown oidclogin provider")
	}

	log.Println("Found oidclogin provider information")

	ctxReady, ready := context.WithCancel(context.Background())

	server := oidclogin.NewServer(cmd.Callback)

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error {
		log.Println("Starting webserver for logout callback")

		resp, err := server.Run(context.TODO(), ready)
		if err != nil {
			fmt.Println("Failed to start server:", err)
		}

		log.Println("Callback received")

		if resp.Error != "" {
			return fmt.Errorf("logout failed: error_code=%s error_description=%s", resp.Code, resp.ErrorDescription)
		}

		log.Println("Successfully logged out")

		return nil
	})

	group.Go(func() error {
		<-ctxReady.Done()

		log.Println("Opening browser session")

		logoutURL, err := url.Parse(providerInfo.Cognito.LogoutURL)
		if err != nil {
			return fmt.Errorf("failed to parse logout URL %w", err)
		}

		queryParams := url.Values{}
		queryParams.Set("client_id", providerInfo.Cognito.ClientID)
		queryParams.Set("logout_uri", cmd.Callback)

		logoutURL.RawQuery = queryParams.Encode()

		fmt.Println("Opening logout URL in your browser:", logoutURL)

		return open.Run(logoutURL.String())
	})

	err = group.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for logout: %w", err)
	}

	return nil
}
