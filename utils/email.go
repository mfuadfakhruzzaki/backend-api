package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// letters berisi karakter yang akan digunakan untuk menghasilkan kode verifikasi
const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateVerificationCode menghasilkan kode verifikasi 6 karakter alfanumerik
func GenerateVerificationCode() string {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // Membuat generator acak lokal
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[randGen.Intn(len(letters))] // Pilih karakter acak dari letters
	}
	return string(code)
}

// SendVerificationEmail mengirimkan email verifikasi dengan kode ke pengguna
func SendVerificationEmail(recipientEmail string, verificationCode string) error {
	// Mendapatkan konfigurasi SMTP dari environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_SENDER")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	// Verifikasi apakah variabel environment telah diatur dengan benar
	if smtpHost == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		log.Fatal("SMTP configuration is missing in environment variables")
		return fmt.Errorf("SMTP configuration is missing in environment variables")
	}

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Email Verification for Data Quota Tracker")
	m.SetBody("text/plain", fmt.Sprintf("Welcome to Data Quota Tracker!\n\nYour verification code is: %s\n\nPlease enter this code to verify your email and start using the app.", verificationCode))

	// Membuat dialer untuk mengirim email
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Invalid SMTP port: %v", err)
		return fmt.Errorf("invalid SMTP port: %v", err)
	}
	d := gomail.NewDialer(smtpHost, port, senderEmail, senderPassword)

	// Kirim email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send verification email to %s: %v", recipientEmail, err)
		return err
	}

	fmt.Printf("Verification email sent to %s with code %s\n", recipientEmail, verificationCode)
	return nil
}
