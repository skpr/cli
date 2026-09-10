package logout

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/skpr/api/pb"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/sync/errgroup"

	oidclogin "github.com/skpr/cli/internal/aws/cognito/oidc/login"
	"github.com/skpr/cli/internal/client"
	"github.com/skpr/cli/internal/client/config"
	credentialscache "github.com/skpr/cli/internal/client/credentials/cache"
)

// CallbackTimeout for how long we wait for the logout callback.
const CallbackTimeout = 30 * time.Second

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

	credentials, found, err := credentialscache.Get(config.API.Host())
	if err != nil {
		return fmt.Errorf("failed to get cached credentials: %w", err)
	}

	if found {
		log.Println("Deleting cached credentials")

		// Signing out with the identity provider is best effort. If our refresh
		// token has already expired, or was revoked, we cannot get an access
		// token to sign out with. We still need to remove the local credentials
		// so that a developer can login again.
		err = globalSignOut(ctx, credentials)
		if err != nil {
			// @todo, We can add this later once we have log levels eg. Turn of debug mode.
			// https://previousnext.atlassian.net/browse/SKPR-1002
			log.Println("Failed to execute global sign out:", err)
		}

		// Delete the file now that we have invalidated our tokens.
		err = credentialscache.Delete(config.API.Host())
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
	defer ready()

	server := oidclogin.NewServer(cmd.Callback)

	// This server is handling a logout callback.
	server.Logout = true

	// The identity provider redirects straight back to us, so unlike login
	// there is no developer interaction to wait for.
	server.CallbackTimeout = CallbackTimeout

	group, groupCtx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		log.Println("Starting webserver for logout callback")

		resp, err := server.Run(context.TODO(), ready)
		if err != nil {
			// The identity provider will not redirect back to us unless our
			// callback URL has been registered as a sign out URL. Our local
			// credentials have already been removed at this point, so we tell
			// the developer what is left over instead of failing.
			if errors.Is(err, oidclogin.ErrCallbackTimeout) {
				log.Println("Did not receive the logout callback:", err)
				log.Println("Local credentials have been removed, however your identity provider browser session may still be active")

				return nil
			}

			return fmt.Errorf("failed to run logout callback server: %w", err)
		}

		log.Println("Callback received")

		if resp.Error != "" {
			return fmt.Errorf("logout failed: error_code=%s error_description=%s", resp.Code, resp.ErrorDescription)
		}

		log.Println("Successfully logged out")

		return nil
	})

	group.Go(func() error {
		// Wait for the callback server to become ready. We also watch the group
		// so that a server which never becomes ready fails the command instead
		// of waiting here forever.
		select {
		case <-ctxReady.Done():
		case <-groupCtx.Done():
			return groupCtx.Err()
		}

		log.Println("Opening browser session")

		logoutURL, err := url.Parse(providerInfo.Cognito.LogoutURL)
		if err != nil {
			return fmt.Errorf("failed to parse logout URL %w", err)
		}

		// Keep any query parameters which the platform has already provided.
		queryParams := logoutURL.Query()
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

// Helper function to sign the developer out of all of their sessions with the
// identity provider.
func globalSignOut(ctx context.Context, credentials credentialscache.Credentials) error {
	token, err := credentials.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(credentials.Cognito.Region), awsconfig.WithCredentialsProvider(aws.AnonymousCredentials{}))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	_, err = cognitoidentityprovider.NewFromConfig(cfg).GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(token.AccessToken),
	})
	if err != nil {
		return fmt.Errorf("failed to global sign out: %w", err)
	}

	return nil
}
