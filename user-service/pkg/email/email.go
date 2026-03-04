package email

import (
	"fmt"

	"github.com/leminhthai/train-ticket/user-service/global"
	"gopkg.in/gomail.v2"
)

func SendOTPEmail(toEmail, otp string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", global.Config.Email.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Mã OTP đặt lại mật khẩu")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Mã OTP của bạn</h2>
		<p>Mã OTP để đặt lại mật khẩu của bạn là:</p>
		<h1 style="color: #4A90E2; letter-spacing: 5px;">%s</h1>
		<p>Mã có hiệu lực trong <strong>5 phút</strong>.</p>
		<p>Nếu bạn không yêu cầu đặt lại mật khẩu, hãy bỏ qua email này.</p>
	`, otp))

	d := gomail.NewDialer(
		global.Config.Email.Host,
		global.Config.Email.Port,
		global.Config.Email.Username,
		global.Config.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}