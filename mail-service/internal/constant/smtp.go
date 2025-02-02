package constant

const (
	SendVerificationSubject = "[Go Microservices] Email Verification"
	AccountVerifiedSubject  = "[Go Microservices] Account Verified"
)

const (
	SendVerificationTemplate = `
<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Email Verification</h1>
		<p>Dear User,</p>
		<p>Thank you for registering. Please use the following OTP code to verify your email address:</p>
		<h2>%v</h2>
		<p>This code will expire in 10 minutes.</p>
		<p>If you did not request this, please ignore this email.</p>
		<p>Best regards,</p>
		<p>The Go Microservices Team</p>
	</body>
</html>
	`

	AccountVerifiedTemplate = `
<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Account Verified</h1>
		<p>Dear User,</p>
		<p>We are pleased to inform you that your account has been successfully verified.</p>
		<p>You can now access all the features of our service.</p>
		<p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>
		<p>Best regards,</p>
		<p>The Go Microservices Team</p>
	</body>
</html>
	`
)
