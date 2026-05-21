package brevo

import "fmt"

func activationCodeHTML(code string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f4f5;font-family:Arial,sans-serif">
  <table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 0">
    <tr><td align="center">
      <table width="480" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,.08)">
        <tr><td style="background:#111827;padding:24px 32px">
          <p style="margin:0;color:#ffffff;font-size:18px;font-weight:bold">Car Rental</p>
        </td></tr>
        <tr><td style="padding:32px">
          <p style="margin:0 0 8px;font-size:22px;font-weight:bold;color:#111827">Verify your email</p>
          <p style="margin:0 0 24px;font-size:14px;color:#6b7280">Use the code below to complete your registration.</p>
          <div style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:6px;padding:20px;text-align:center;margin-bottom:24px">
            <span style="font-size:32px;font-weight:bold;letter-spacing:8px;color:#111827">%s</span>
          </div>
          <p style="margin:0;font-size:13px;color:#9ca3af">This code expires shortly. If you did not request this, you can safely ignore this email.</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, code)
}
