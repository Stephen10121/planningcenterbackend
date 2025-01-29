package email

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/resend/resend-go/v2"
)

func SendMail(toEmail string, subject string, html string) (*resend.SendEmailResponse, error) {
	ctx := context.TODO()
	client := resend.NewClient(initializers.ResendKey)

	params := &resend.SendEmailRequest{
		From:    "EventLooker <planningcenter@mail.stephengruzin.dev>",
		To:      []string{toEmail},
		Subject: subject,
		Html:    html,
	}

	return client.Emails.SendWithContext(ctx, params)
}

// This function sends a warning email to the recipient that their refresh token is about to expire and that they need to login to their service to set a new token.
func TokenExpireWarning(daysLeft int, recipient string) (*resend.SendEmailResponse, error) {
	file, err := os.Open("./email/templates/tokenExpire.html")
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return SendMail(
		recipient,
		"WARNING!! Your refresh token is about to expire in "+strconv.Itoa(daysLeft)+" days.",
		string(b),
	)
}
