package service

import (
	"crypto/tls"
	"fmt"
	"log/slog"

	"gopkg.in/gomail.v2"

	"vniizht/internal/config"
	"vniizht/internal/model"
)

type EmailService struct {
	dialer  *gomail.Dialer
	from    string
	appURL  string
	enabled bool
}

func NewEmailService(cfg *config.Config) *EmailService {
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	// порт 465 — implicit TLS; порт 587/25 — STARTTLS (SSL: false, gomail сам делает STARTTLS)
	if cfg.SMTPPort == 465 {
		d.SSL = true
	} else {
		d.SSL = false
		d.TLSConfig = &tls.Config{ServerName: cfg.SMTPHost}
	}
	return &EmailService{
		dialer:  d,
		from:    cfg.SMTPFrom,
		appURL:  cfg.AppURL,
		enabled: cfg.SMTPEnabled,
	}
}

func (e *EmailService) SendNotification(to string, typ model.NotificationType, payload model.NotificationPayload) {
	if !e.enabled || to == "" {
		return
	}
	subject, body := e.buildEmail(typ, payload)
	if subject == "" {
		return
	}
	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", e.from)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
		if err := e.dialer.DialAndSend(m); err != nil {
			slog.Warn("failed to send email", "to", to, "type", typ, "err", err)
		}
	}()
}

func (e *EmailService) buildEmail(typ model.NotificationType, p model.NotificationPayload) (subject, body string) {
	questionURL := fmt.Sprintf("%s/questions/%d", e.appURL, p.QuestionID)
	settingsURL := e.appURL + "/settings"

	switch typ {
	case model.NotifNewAnswer:
		subject = fmt.Sprintf("Новый ответ на вопрос «%s»", p.QuestionTitle)
		body = e.emailHTML(subject,
			fmt.Sprintf("Пользователь <b>%s</b> ответил на ваш вопрос.", p.ActorUsername),
			fmt.Sprintf("«%s»", p.QuestionTitle),
			questionURL, settingsURL,
		)
	case model.NotifAnswerVerified:
		subject = "Ваш ответ подтверждён как достоверный"
		body = e.emailHTML(subject,
			fmt.Sprintf("Специалист <b>%s</b> отметил ваш ответ как достоверный.", p.ActorUsername),
			fmt.Sprintf("«%s»", p.QuestionTitle),
			questionURL, settingsURL,
		)
	case model.NotifQuestionAssigned:
		subject = "Вам назначен новый вопрос"
		body = e.emailHTML(subject,
			fmt.Sprintf("Пользователь <b>%s</b> задал вопрос, назначенный вам.", p.ActorUsername),
			fmt.Sprintf("«%s»", p.QuestionTitle),
			questionURL, settingsURL,
		)
	case model.NotifQuestionClosed:
		subject = "Ваш вопрос закрыт"
		body = e.emailHTML(subject,
			fmt.Sprintf("Специалист <b>%s</b> закрыл ваш вопрос.", p.ActorUsername),
			fmt.Sprintf("«%s»", p.QuestionTitle),
			questionURL, settingsURL,
		)
	case model.NotifNewComment:
		subject = "Новый комментарий к вашему ответу"
		body = e.emailHTML(subject,
			fmt.Sprintf("Пользователь <b>%s</b> оставил комментарий к вашему ответу.", p.ActorUsername),
			fmt.Sprintf("«%s»", p.QuestionTitle),
			questionURL, settingsURL,
		)
	}
	return
}

func (e *EmailService) emailHTML(title, lead, subtitle, questionURL, settingsURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;background:#f5f5f5;padding:24px">
  <div style="max-width:560px;margin:0 auto;background:#fff;border-radius:8px;padding:32px;border:1px solid #e0e0e0">
    <div style="margin-bottom:16px">
      <span style="background:#2563eb;color:#fff;font-weight:bold;padding:4px 10px;border-radius:4px;font-size:13px">ИСС АПК «ЭЛЬБРУС»</span>
    </div>
    <h2 style="margin:0 0 12px;font-size:18px;color:#111">%s</h2>
    <p style="margin:0 0 8px;color:#444;font-size:15px">%s</p>
    <p style="margin:0 0 24px;color:#666;font-size:14px">%s</p>
    <a href="%s" style="display:inline-block;background:#2563eb;color:#fff;text-decoration:none;padding:10px 20px;border-radius:6px;font-size:14px;font-weight:bold">Перейти к вопросу</a>
    <p style="margin:24px 0 0;color:#999;font-size:12px">
      Вы получили это письмо, потому что подписаны на уведомления.<br>
      Изменить настройки можно в <a href="%s" style="color:#2563eb">настройках профиля</a>.
    </p>
  </div>
</body>
</html>`, title, lead, subtitle, questionURL, settingsURL)
}
