package optimizer

import (
	"strings"
	"unicode"
)

//TODO rewrite css parser to optimize code size and include this tokenizer inside as method for for Selector

// skipPseudo skips over a pseudo-class or pseudo-element in the given string.
// It assumes that s[start] is ':' and returns the index of the last character
// that belongs to the pseudo selector.
func skipPseudo(s string, start int) int {
	n := len(s)
	i := start
	if i < n && s[i] == ':' {
		i++
		// If it's a pseudo-element with two colons, skip the second.
		if i < n && s[i] == ':' {
			i++
		}
		// Skip the pseudo name: letters, digits, hyphen or underscore.
		for i < n {
			r := rune(s[i])
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
				i++
			} else if i < n && s[i] == '(' {
				// If it's a function-like pseudo-class (e.g. :nth-child(2)),
				// skip until the matching ')'.
				i++ // skip '('
				depth := 1
				for i < n && depth > 0 {
					if s[i] == '(' {
						depth++
					} else if s[i] == ')' {
						depth--
					}
					i++
				}
			} else {
				break
			}
		}
	}
	// Return the index of the last character that is part of the pseudo.
	return i - 1
}

// splitCompoundSelector splits a compound selector (a sequence of simple selectors
// with no combinators) into individual tokens. Pseudo selectors are removed.
// For example, "div#main.content[data-type=\"example\"]" becomes
// ["div", "#main", ".content[data-type=\"example\"]"].
func splitCompoundSelector(compound string) []string {
	var tokens []string
	var current strings.Builder
	n := len(compound)
	i := 0
	for i < n {
		r := rune(compound[i])
		if r == ':' {
			// Before skipping the pseudo, flush the current token if any.
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			// Skip the pseudo selector.
			newIndex := skipPseudo(compound, i)
			i = newIndex + 1
			continue
		}
		// A new simple selector begins at '#' or '.'.
		if r == '#' || r == '.' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			current.WriteRune(r)
		} else if r == '[' {
			// Append the attribute selector to the current token.
			current.WriteRune(r)
			i++
			// Append until the matching ']' is found.
			for i < n && compound[i] != ']' {
				current.WriteByte(compound[i])
				i++
			}
			if i < n && compound[i] == ']' {
				current.WriteByte(compound[i])
			}
		} else {
			// Accumulate any other character (part of an element name, etc.)
			current.WriteRune(r)
		}
		i++
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

// tokenizeCSSSelector splits a CSS selector into tokens.
// It separates compound segments (which are further split into simple selectors)
// from combinators (whitespace, ">", "+", "~", ",") and removes pseudo selectors.
func tokenizeCSSSelector(selector string) []string {
	var finalTokens []string
	var compound strings.Builder

	// State flags for attribute selectors.
	inAttr := false
	inQuote := false
	var quoteChar rune

	runes := []rune(selector)
	n := len(runes)
	i := 0
	for i < n {
		r := runes[i]
		// Handle attribute selectors: once inside, include everything until the matching ']'.
		if inQuote {
			compound.WriteRune(r)
			if r == quoteChar {
				inQuote = false
			}
			i++
			continue
		}
		if inAttr {
			compound.WriteRune(r)
			if r == '\'' || r == '"' {
				inQuote = true
				quoteChar = r
			} else if r == ']' {
				inAttr = false
			}
			i++
			continue
		}
		if r == '[' {
			inAttr = true
			compound.WriteRune(r)
			i++
			continue
		}
		// When a pseudo selector is encountered, flush the current compound
		// and skip the pseudo.
		if r == ':' {
			if compound.Len() > 0 {
				compStr := compound.String()
				simpleSelectors := splitCompoundSelector(compStr)
				finalTokens = append(finalTokens, simpleSelectors...)
				compound.Reset()
			}
			// Skip the pseudo selector.
			remaining := string(runes[i:])
			newIndex := skipPseudo(remaining, 0)
			i += newIndex + 1
			continue
		}
		// If we hit a combinator (">", "+", "~", ",") or whitespace, then flush the compound.
		if r == '>' || r == '+' || r == '~' || r == ',' || unicode.IsSpace(r) {
			if compound.Len() > 0 {
				compStr := compound.String()
				simpleSelectors := splitCompoundSelector(compStr)
				finalTokens = append(finalTokens, simpleSelectors...)
				compound.Reset()
			}
			// For whitespace, add a single space token (avoid duplicates).
			if !unicode.IsSpace(r) {
				finalTokens = append(finalTokens, string(r))
			}
			i++
			continue
		}
		// Otherwise, accumulate the character into the current compound selector.
		compound.WriteRune(r)
		i++
	}
	// Flush any remaining compound.
	if compound.Len() > 0 {
		compStr := compound.String()
		simpleSelectors := splitCompoundSelector(compStr)
		finalTokens = append(finalTokens, simpleSelectors...)
	}
	return finalTokens
}
