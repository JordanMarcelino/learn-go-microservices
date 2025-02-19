package constant

const (
	SendVerificationSubject = "[Go Microservices] Email Verification"
	AccountVerifiedSubject  = "[Go Microservices] Account Verified"
	PaymentReminderSubject  = "[Go Microservices] Payment Reminder"
	OrderCancelledSubject   = "[Go Microservices] Order Cancelled"
	OrderSuccessSubject     = "[Go Microservices] Order Success"
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

	PaymentReminderTemplate = `
<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Payment Reminder</h1>
		<p>Dear User,</p>
		<p>This is a friendly reminder that your payment is due. Below are the details of your payment:</p>
		<p>Order ID: %v</p>
		<p>Total Amount: %v</p>
		<p>Due Date: %v</p>
		<p>Please make the payment by the due date to avoid any late fees.</p>
		<p>If you have already made the payment, please disregard this email.</p>
		<p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>
		<p>Best regards,</p>
		<p>The Go Microservices Team</p>
	</body>
</html>
	`

	OrderCancelledTemplate = `
<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Order Cancelled</h1>
		<p>Dear User,</p>
		<p>We regret to inform you that your order has been cancelled. Below are the details of your order:</p>
		<p>Order ID: %v</p>
		<p>Total Amount: %v</p>
		<p>Description: %v</p>
		<p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>
		<p>Best regards,</p>
		<p>The Go Microservices Team</p>
	</body>
</html>
	`

	OrderSuccessTemplate = `
<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Order Success</h1>
		<p>Dear User,</p>
		<p>We are pleased to inform you that your order has been successfully placed. Below are the details of your order:</p>
		<p>Order ID: %v</p>
		<p>Total Amount: %v</p>
		<p>Description: %v</p>
		<p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>
		<p>Best regards,</p>
		<p>The Go Microservices Team</p>
	</body>
</html>
	`
)
