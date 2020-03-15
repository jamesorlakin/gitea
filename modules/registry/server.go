package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"code.gitea.io/gitea/modules/log"
)

// AuthServer is the token authentication server
type AuthServer struct {
	authorizer     Authorizer
	authenticator  Authenticator
	tokenGenerator TokenGenerator
	crt, key       string
}

// NewAuthServer creates a new AuthServer
func NewAuthServer(opt *Option) (*AuthServer, error) {
	pb, prk, err := loadCertAndKey(opt.Certfile, opt.Keyfile)
	if err != nil {
		return nil, err
	}
	tk := &TokenOption{Expire: opt.TokenExpiration, Issuer: opt.TokenIssuer}
	generator := newTokenGenerator(pb, prk, tk)
	return &AuthServer{
		authorizer:     &GiteaAuthorizer{},
		authenticator:  &GiteaAuthenticator{},
		tokenGenerator: generator, crt: opt.Certfile, key: opt.Keyfile,
	}, nil
}

func (srv *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := srv.parseRequest(r)
	// grab user's auth parameters
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "registry: please provide Gitea credentials to interface with this registry", http.StatusUnauthorized)
		return
	}
	user, err := srv.authenticator.Authenticate(username, password)
	if err != nil {
		http.Error(w, "registry: invalid Gitea credentials", http.StatusUnauthorized)
		return
	}
	actions, err := srv.authorizer.Authorize(user, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// create token for this user using the actions returned
	// from the authorization check
	tk, err := srv.tokenGenerator.Generate(req, actions)
	if err != nil {
		log.Error("registry: couldn't generate token: %s", err.Error())
		http.Error(w, "registry: internal server error", http.StatusInternalServerError)
		return
	}
	srv.ok(w, tk)
}

func (srv *AuthServer) parseRequest(r *http.Request) *AuthorizationRequest {
	q := r.URL.Query()
	req := &AuthorizationRequest{
		Service: q.Get("service"),
		Account: q.Get("account"),
		IP:      r.RemoteAddr,
	}
	parts := strings.Split(r.URL.Query().Get("scope"), ":")
	if len(parts) > 0 {
		req.Type = parts[0]
	}
	if len(parts) > 1 {
		req.Name = parts[1]
	}
	if len(parts) > 2 {
		req.Actions = strings.Split(parts[2], ",")
	}
	if req.Account == "" {
		req.Account = req.Name
	}
	return req
}

func (srv *AuthServer) Run(addr string) error {
	fmt.Printf("registry: authentication server running at %s", addr)
	return http.ListenAndServe(addr, srv)
	// return http.ListenAndServeTLS(addr, srv.crt, srv.key, nil)
}

func (srv *AuthServer) ok(w http.ResponseWriter, tk *Token) {
	data, _ := json.Marshal(tk)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func encodeBase64(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
