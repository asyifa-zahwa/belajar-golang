package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unioffice/document"

	//"github.com/unidoc/unipdf/document"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// Pastikan untuk mengganti dengan lisensi UniPDF Anda
const uniPDFLicense = "YOUR_UNIPDF_LICENSE_KEY"

func ConvertToDocx() {
	// Inisialisasi lisensi UniPDF
	err := license.SetLicenseKey(uniPDFLicense, "")
	if err != nil {
		log.Fatal(err)
	}

	// Membaca file PDF
	filePath := "input.pdf"
	pdfFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open PDF file: %v", err)
	}
	defer pdfFile.Close()

	// Ekstraksi teks dari PDF
	pdfReader, err := model.NewPdfReader(pdfFile)
	if err != nil {
		log.Fatalf("Unable to create PDF reader: %v", err)
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		fmt.Println("Error getting number of pages:", err)
		return
	}

	// Membaca seluruh teks dari PDF
	var extractedText string
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			log.Fatalf("Unable to get page %d: %v", pageNum, err)
		}

		extractor, err := extractor.New(page)
		if err != nil {
			log.Fatalf("i dont know but error", err)
		}
		text, err := extractor.ExtractText()
		if err != nil {
			log.Fatalf("Unable to extract text from page %d: %v", pageNum, err)
		}
		extractedText += text
	}

	// Buat file DOCX baru
	doc := document.New()

	// Tambahkan teks yang diekstraksi ke dalam dokumen DOCX
	doc.AddParagraph().AddRun().AddText(extractedText)

	// Simpan DOCX ke file
	outputPath := "output.docx"
	err = doc.SaveToFile(outputPath)
	if err != nil {
		log.Fatalf("Unable to save DOCX file: %v", err)
	}

	fmt.Println("PDF berhasil dikonversi ke DOCX:", outputPath)
}
