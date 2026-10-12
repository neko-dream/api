package delivery

import (
	"fmt"
	"net"
	"strconv"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/email"
	"github.com/neko-dream/api/internal/infrastructure/external/aws/ses"
)

// NewSender keeps SES in deployed environments; local development must opt
// into a loopback SMTP receiver rather than sending real email by accident.
func NewSender(cfg *config.Config) (email.EmailSender, error) {
	if cfg == nil {
		return nil, fmt.Errorf("email configuration is missing")
	}
	if cfg.Env != config.LOCAL {
		return ses.NewSESEmailSender(cfg)
	}
	host, port, err := net.SplitHostPort(cfg.SMTP_ADDR)
	ip := net.ParseIP(host)
	portNumber, portErr := strconv.Atoi(port)
	if err != nil || portErr != nil || portNumber < 1 || portNumber > 65535 || (host != "localhost" && (ip == nil || !ip.IsLoopback())) {
		return nil, fmt.Errorf("local SMTP_ADDR must specify a loopback host and port")
	}
	return &smtpSender{cfg: cfg}, nil
}
