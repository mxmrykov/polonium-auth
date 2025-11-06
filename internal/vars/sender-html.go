package vars

const (
	VerificationCode = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Arial, sans-serif; background-color: #f6f9fc;">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f6f9fc; padding: 50px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); overflow: hidden;">
                    <!-- Header -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 40px 0; text-align: center;">
                            <h1 style="color: #000000; margin: 0; font-size: 28px; font-weight: 600;">Verify Your Email</h1>
                            <p style="color: #888888; margin: 10px 0 0 0; font-size: 16px;">Complete your registration</p>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <p style="color: #333333; font-size: 16px; line-height: 1.6; margin: 0 0 20px 0;">
                                Hello!
                            </p>
                            <p style="color: #555555; font-size: 16px; line-height: 1.6; margin: 0 0 25px 0;">
                                Thank you for registering! To complete your verification, please use the following 6-digit verification code:
                            </p>
                            
                            <!-- Verification Code Box -->
                            <table width="100%%" cellpadding="0" cellspacing="0" border="0">
                                <tr>
                                    <td align="center">
                                        <div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #000000; font-size: 32px; font-weight: bold; letter-spacing: 8px; padding: 20px 30px; border-radius: 8px; display: inline-block; margin: 20px 0;">
                                            %s
                                        </div>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #888888; font-size: 14px; line-height: 1.5; margin: 25px 0 0 0;">
                                This code will expire in 5 minutes. If you didn't request this verification, please ignore this email.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8f9fa; padding: 25px 30px; border-top: 1px solid #eaeaea;">
                            <p style="color: #999999; font-size: 12px; line-height: 1.4; margin: 0; text-align: center;">
                                &copy; 2025 Polonium. All rights reserved.<br>
                                If you have any questions, contact us at support@polonium.ws
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`
)
