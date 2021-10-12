package context

import (
	"emviwiki/backend/client"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"strings"
)

// EmviContext contains contextual information for all use cases.
type EmviContext struct {
	Organization  *model.Organization
	UserId        hide.ID
	Scopes        map[string]client.Scope
	TrustedClient bool
}

// NewEmviContext returns a new context for given parameters.
// The scopes are converted to a valid client.Scope map with their name as index
// and must be in the format "name:rw". Invalid scopes are ignored.
func NewEmviContext(orga *model.Organization, userId hide.ID, scopes []string, trustedClient bool) EmviContext {
	return EmviContext{orga, userId, toScopeMap(scopes), trustedClient}
}

func NewEmviUserContext(orga *model.Organization, userId hide.ID) EmviContext {
	return EmviContext{Organization: orga, UserId: userId}
}

// IsUser returns true if this context was initiated by a user request.
func (ctx *EmviContext) IsUser() bool {
	return ctx.UserId != 0
}

// IsClient returns true if this context was initiated by a client request.
func (ctx *EmviContext) IsClient() bool {
	return ctx.UserId == 0
}

// HasScopes checks if the client of this context has given scopes.
// This checks for exact matches, so read and write permissions on scopes must match.
// If this context was created by a user request or the client is trusted true is returned.
// If this context was created by a client and no scopes are passed or it's an entry organization false is returned.
func (ctx *EmviContext) HasScopes(scopes ...client.Scope) bool {
	if ctx.UserId != 0 || ctx.TrustedClient {
		return true
	} else if !ctx.Organization.Expert {
		return false
	}

	found := 0

	for _, scope := range scopes {
		result, ok := ctx.Scopes[scope.Name]

		if ok && result.Read == scope.Read && result.Write == scope.Write {
			found++

			if found == len(scopes) {
				return true
			}
		}
	}

	return false
}

func toScopeMap(scopes []string) map[string]client.Scope {
	scopeMap := make(map[string]client.Scope)

	for _, scope := range scopes {
		parts := strings.Split(scope, ":")

		if len(parts) == 2 {
			_, ok := client.Scopes[parts[0]]

			if ok {
				if parts[1] == "r" {
					scopeMap[parts[0]] = client.Scope{parts[0], true, false}
				} else if parts[1] == "rw" {
					scopeMap[parts[0]] = client.Scope{parts[0], true, true}
				}
			}
		}
	}

	return scopeMap
}
