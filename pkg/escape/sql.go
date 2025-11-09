package escape

import "strings"

var likeEscaper = strings.NewReplacer(
	`\`, `\\`,
	`%`, `\%`,
	`_`, `\_`,
)

func EscapeLikePattern(s string) string {
	return likeEscaper.Replace(s)
}
