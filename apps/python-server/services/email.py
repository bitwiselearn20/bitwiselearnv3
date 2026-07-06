import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from config import get_settings

settings = get_settings()


def send_email(to_email: str, subject: str, html_body: str):
    msg = MIMEMultipart("alternative")
    msg["Subject"] = subject
    msg["From"] = settings.EMAIL_USER
    msg["To"] = to_email
    msg.attach(MIMEText(html_body, "html"))

    with smtplib.SMTP_SSL("smtp.gmail.com", 465) as server:
        server.login(settings.EMAIL_USER, settings.EMAIL_PASS)
        server.sendmail(settings.EMAIL_USER, to_email, msg.as_string())


def send_welcome_email(to_email: str, name: str, password: str, role: str):
    html = f"""
    <h2>Welcome to BitwiseLearn!</h2>
    <p>Hello {name},</p>
    <p>Your {role} account has been created successfully.</p>
    <p><strong>Email:</strong> {to_email}</p>
    <p><strong>Password:</strong> {password}</p>
    <p>Please change your password after your first login.</p>
    <p>Best regards,<br>BitwiseLearn Team</p>
    """
    send_email(to_email, f"Welcome to BitwiseLearn - {role} Account", html)


def send_otp_email(to_email: str, otp: str):
    html = f"""
    <h2>Password Reset OTP</h2>
    <p>Your OTP for password reset is:</p>
    <h1 style="color: #4F46E5; letter-spacing: 4px;">{otp}</h1>
    <p>This OTP is valid for 10 minutes.</p>
    <p>If you didn't request this, please ignore this email.</p>
    <p>Best regards,<br>BitwiseLearn Team</p>
    """
    send_email(to_email, "BitwiseLearn - Password Reset OTP", html)


def send_contact_email(name: str, email: str, message: str):
    html = f"""
    <h2>New Contact Form Submission</h2>
    <p><strong>Name:</strong> {name}</p>
    <p><strong>Email:</strong> {email}</p>
    <p><strong>Message:</strong> {message}</p>
    """
    send_email(settings.EMAIL_USER, f"Contact Form: {name}", html)
