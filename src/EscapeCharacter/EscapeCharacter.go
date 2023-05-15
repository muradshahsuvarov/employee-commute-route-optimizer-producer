package EscapeCharacter

import "strings"

// Helper function to escape special characters in the message
func EscapeSpecialCharacters(message string) string {
	escapedMessage := strings.ReplaceAll(message, `\`, `\\`)
	escapedMessage = strings.ReplaceAll(escapedMessage, `"`, `\"`)
	escapedMessage = strings.ReplaceAll(escapedMessage, `'`, `\'`)
	escapedMessage = strings.ReplaceAll(escapedMessage, "\n", `\n`)
	escapedMessage = strings.ReplaceAll(escapedMessage, "\r", `\r`)
	escapedMessage = strings.ReplaceAll(escapedMessage, "\t", `\t`)
	return escapedMessage
}
