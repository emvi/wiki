package client

import "strings"

var (
	// list of valid scopes and if read/write is supported
	Scopes = map[string]Scope{
		"organization":          {"organization", true, false},
		"language":              {"language", true, false},
		"articles":              {"articles", true, false},
		"article_authors":       {"article_authors", true, false},
		"article_authors_mails": {"article_authors_mails", true, false},
		"article_history":       {"article_history", true, false},
		"lists":                 {"lists", true, false},
		"tags":                  {"tags", true, false},
		"pinned":                {"pinned_articles", true, false},
		"search_articles":       {"search_articles", true, false},
		"search_lists":          {"search_lists", true, false},
		"search_tags":           {"search_tags", true, false},
		"search_all":            {"search_all", true, false},
	}
)

// Scope is a client scope.
type Scope struct {
	Name  string `json:"name"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}

// String returns the string representation of this scope in the format "name:rw".
func (scope Scope) String() string {
	if scope.Write {
		return scope.Name + ":rw"
	} else if scope.Read {
		return scope.Name + ":r"
	}

	return scope.Name
}

// ScopeFromString converts a string to Scope.
// The scope must be represented as "name:rw".
// If it cannot be converted an empty Scope object is returned.
func ScopeFromString(scopeStr string) Scope {
	scope := Scope{}
	parts := strings.Split(scopeStr, ":")

	if len(parts) == 2 {
		_, ok := Scopes[parts[0]]

		if ok {
			scope.Name = parts[0]

			if parts[1] == "rw" {
				scope.Read = true
				scope.Write = true
			} else if parts[1] == "r" {
				scope.Read = true
			}
		}
	}

	return scope
}
