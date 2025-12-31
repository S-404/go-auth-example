package validation

const (
	ValidPassword = "valid_password"
	ValidUsername = "valid_username"
)

var UsernameValidation = validateRegex("^[a-zA-Z0-9]{8,20}$")
var PasswordValidation = validateRegexPatterns([]string{".{8,32}", "[a-z]", "[!@#$%^&]", "[A-Z]", "[0-9]", "[^\\d\\w]"})
