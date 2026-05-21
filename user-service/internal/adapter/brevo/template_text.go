package brevo

import "fmt"

func activationCodeText(code string) string {
	return fmt.Sprintf(
		"Car Rental — Email Verification\n\n"+
			"Your activation code is:\n\n"+
			"  %s\n\n"+
			"Enter this code to verify your email address. It expires shortly.\n\n"+
			"If you did not request this, you can safely ignore this message.",
		code,
	)
}
