package email

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gustavoz65/finext/internal/config"
	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	zerolog "github.com/rs/zerolog"
)

type Client struct {
	client        *resend.Client
	logger        *zerolog.Logger
	senderName    string
	senderEmail   string
	templatesPath string
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	return &Client{
		client:        resend.NewClient(cfg.Integration.ResendAPIKey),
		logger:        logger,
		senderName:    cfg.Integration.SenderName,
		senderEmail:   cfg.Integration.SenderEmail,
		templatesPath: cfg.Integration.TemplatesPath,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s.html", c.templatesPath, templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "Falied to parse email template %s", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", c.senderName, c.senderEmail),
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	_, err = c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendWelcomeEmail envia email de boas-vindas para o usuario
func (c *Client) SendWelcomeEmail(to, firstName string) error {
	data := map[string]string{
		"FirstName": firstName,
	}
	return c.SendEmail(to, "Bem-vindo ao Cashing!", TemplateWelcome, data)
}
