// internal/domain/services/pdfgen/certificate.go

package pdfgen

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gitlab.com/w0ikid/study-platform/internal/domain/models"
)

func GenerateCertificatePDF(user *models.User, course *models.Course, issuedAt time.Time) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Background
	templatePath := "/app/internal/domain/services/pdfgen/background.jpg"
	// templatePath := "/home/doni/KBTU/sem4/golangKBTU/lol/study-platform/internal/domain/services/pdfgen/background.jpg"
	pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")

	// Title
	pdf.SetFont("Times", "B", 40)
	pdf.SetTextColor(0, 102, 204)
	pdf.SetXY(0, 50)
	pdf.CellFormat(297, 20, "Certificate of Achievement", "0", 1, "C", false, 0, "")

	// Subtitle
	pdf.SetFont("Arial", "", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(0, 80)
	pdf.CellFormat(297, 10, "This is to certify that", "0", 1, "C", false, 0, "")

	// Name
	pdf.SetFont("Arial", "B", 35)
	pdf.SetTextColor(255, 0, 0)
	pdf.SetXY(0, 100)
	pdf.CellFormat(297, 15, user.Name+" "+user.Surname, "0", 1, "C", false, 0, "")

	// Program
	pdf.SetFont("Arial", "", 18)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(0, 120)
	pdf.CellFormat(297, 10, fmt.Sprintf("Has successfully completed the course: %s", course.Name), "0", 1, "C", false, 0, "")

	// Date & Signature
	pdf.SetFont("Arial", "I", 16)
	pdf.SetXY(50, 160)
	pdf.CellFormat(100, 10, fmt.Sprintf("Date: %s", issuedAt.Format("2006-01-02")), "0", 0, "L", false, 0, "")

	pdf.SetXY(180, 160)
	pdf.CellFormat(100, 10, "Signature: ____________", "0", 1, "R", false, 0, "")

	log.Printf("PDF generated for user %s %s course %s", user.Name, user.Surname, course.Name)

	// Output to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return nil, err
	}
	log.Printf("PDF generated successfully for user %s %s", user.Name, user.Surname)
	return buf.Bytes(), nil
}
