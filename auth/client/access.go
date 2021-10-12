package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/shared/constants"
	"fmt"
	"github.com/emvi/hide"
)

type ClientAccess struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func ValidateClientCredentials(granttype, clientId, clientSecret string) (*ClientAccess, error) {
	if granttype != constants.AuthGrantType {
		return nil, errs.GrantTypeInvalid
	}

	client := model.GetClientByClientIdAndClientSecret(clientId, clientSecret)

	if client == nil {
		return nil, errs.ClientCredentialsInvalid
	}

	token, ttl, err := jwt.NewClientToken(&jwt.ClientTokenClaims{
		ClientId: clientId,
		Trusted:  client.Trusted,
		Scopes:   getClientScopes(client.ID),
	})

	if err != nil {
		return nil, err
	}

	return &ClientAccess{constants.AuthTokenType,
		ttl.Unix(),
		token}, nil
}

func getClientScopes(id hide.ID) []string {
	scopes := model.FindScopeByClientId(id)
	scopeStrs := make([]string, 0, len(scopes))

	for _, scope := range scopes {
		scopeStrs = append(scopeStrs, fmt.Sprintf("%s:%s", scope.Key, scope.Value))
	}

	return scopeStrs
}
