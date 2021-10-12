package db

const (
	space                    = ' '
	tab                      = '\t'
	tsvectorOr               = "|"
	tsvectorAnd              = "&"
	tsvectorNot              = "!"
	tsvectorPrefix           = "*"
	tsvectorLexeme           = ":"
	tsvectorParenthesesOpen  = "("
	tsvectorParenthesesClose = ")"
)

func ToTSVector(query string) string {
	tokens := parseTSVectorTokens(query)
	result := ""
	wasOperator := false

	for i, token := range tokens {
		isOperator := token == tsvectorOr || token == tsvectorAnd || token == tsvectorNot

		if isOperator && !wasOperator && i != len(tokens)-1 {
			if token == tsvectorNot {
				if len(result) != 0 {
					result += tsvectorAnd + tsvectorNot
				} else {
					result += tsvectorNot
				}
			} else if len(result) != 0 {
				result += token
			}

			wasOperator = true
		} else if !isOperator {
			if !wasOperator && len(result) != 0 {
				result += string(tsvectorOr)
			}

			result += token
			wasOperator = false
		}
	}

	return result
}

func parseTSVectorTokens(query string) []string {
	tokens := make([]string, 0)
	token := ""

	for _, c := range query {
		cStr := string(c)

		if tsvectorRemoveCharacter(cStr) {
			c = space
			cStr = string(space)
		}

		isOperator := cStr == tsvectorOr || cStr == tsvectorAnd || cStr == tsvectorNot
		isDelimiter := c == space || c == tab || isOperator

		if isDelimiter {
			if len(token) != 0 {
				tokens = append(tokens, token)
			}

			if isOperator {
				tokens = append(tokens, cStr)
			}

			token = ""
		} else {
			if c != space && c != tab {
				token += cStr
			}
		}
	}

	// append last token
	if len(token) != 0 {
		tokens = append(tokens, token)
	}

	return tokens
}

func tsvectorRemoveCharacter(c string) bool {
	return c == tsvectorPrefix ||
		c == tsvectorLexeme ||
		c == tsvectorParenthesesOpen ||
		c == tsvectorParenthesesClose
}
