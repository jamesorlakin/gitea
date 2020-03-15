package registry

import (
	"errors"
	"strings"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/log"
)

// StartRegistry kicks off a separate HTTP server for Docker authentication using Gitea's auth functions
func StartRegistry() {
	log.Info("registry: %s", "running")
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return
	// }
	serverOptions := &Option{
		Certfile:        "C:\\Users\\LAKINJ\\go\\src\\code.gitea.io\\gitea\\RootCA.pem",
		Keyfile:         "C:\\Users\\LAKINJ\\go\\src\\code.gitea.io\\gitea\\RootCA.key",
		TokenExpiration: 30,
		TokenIssuer:     "gitea-issuer",
	}
	server, err := NewAuthServer(serverOptions)
	if err != nil {
		return
	}
	server.Run("0.0.0.0:5001")
}

// GiteaAuthenticator makes authentication successful by default
type GiteaAuthenticator struct{}

func (d *GiteaAuthenticator) Authenticate(username string, password string) (*models.User, error) {
	log.Info("registry: attempting authentication for %s", username)
	user, err := models.UserSignIn(username, password)
	if (err) != nil {
		log.Warn("registry: couldn't authenticate %s", username)
		return nil, err
	}
	return user, nil
}

// GiteaAuthorizer makes authorization successful by default
type GiteaAuthorizer struct{}

func (d *GiteaAuthorizer) Authorize(user *models.User, req *AuthorizationRequest) ([]string, error) {
	// If running `docker login` there's no repo info. Skip granting permissions.
	if req.Name == "" {
		return []string{}, nil
	}

	path := strings.Split(req.Name, "/")
	if len(path) < 2 {
		return []string{}, errors.New("registry: image name must be in 'owner/repo' format")
	}
	repo, err := models.GetRepositoryByOwnerAndName(path[0], path[1])
	if err != nil {
		log.Warn("registry: couldn't find repository %s", req.Name)
		return []string{}, err
	}
	perm, err := models.GetUserRepoPermission(repo, user)
	if err != nil {
		log.Warn("registry: couldn't get permissions for %s @ %s", user.LoginName, repo.Name)
		return []string{}, err
	}

	actions := []string{}
	canPull := user.IsAdmin || perm.CanRead(models.UnitTypeCode)
	if canPull {
		actions = append(actions, "pull")
	}

	canPush := user.IsAdmin || perm.CanWrite(models.UnitTypeCode)
	if canPush {
		actions = append(actions, "push")
	}

	return actions, nil
}
