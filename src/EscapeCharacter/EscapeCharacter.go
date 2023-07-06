package EscapeCharacter

import "strings"

// Helper function to escape special characters in the message

func EscapeSpecialCharacters(message string) string {
	escapedMessage := strings.ReplaceAll(message, `"`, `\"`)
	return escapedMessage
}
