package controllers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unidoc/unioffice/document"

	//"github.com/unidoc/unipdf/document"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// Pastikan untuk mengganti dengan lisensi UniPDF Anda
const uniPDFLicense = "84fa0a89996aa28922e938277081f81312bfe7433380b93bdd22229b5467baed"

func ConvertToDocx(folderName string) (string, error) {
	// Inisialisasi lisensi UniPDF
	err := license.SetLicenseKey(uniPDFLicense, "")
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	// Membaca file PDF
	folderPatch := "./output/" + folderName + "/"
	// Path folder yang ingin diperiksa
	//folderPath := "./folderA"

	// Membaca isi folder
	files, err := ioutil.ReadDir(folderPatch)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	// Menampilkan nama file di folder
	namaFile := ""
	fmt.Println("Daftar file dalam folder:")
	for _, file := range files {

		namaFile = file.Name()
		//fmt.Println(file.Name())
		return namaFile, nil
	}

	pdfFile, err := os.Open(namaFile)
	if err != nil {
		// log.Fatalf("Unable to open PDF file: %v", err)
		return "", err
	}
	defer pdfFile.Close()

	// Ekstraksi teks dari PDF
	pdfReader, err := model.NewPdfReader(pdfFile)
	if err != nil {
		// log.Fatalf("Unable to create PDF reader: %v", err)
		return "", err
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		// fmt.Println("Error getting number of pages:", err)
		return "", err
	}

	// Membaca seluruh teks dari PDF
	var extractedText string
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			// log.Fatalf("Unable to get page %d: %v", pageNum, err)
			return "", err
		}

		extractor, err := extractor.New(page)
		if err != nil {
			// log.Fatalf("i dont know but error", err)
			return "", err
		}
		text, err := extractor.ExtractText()
		if err != nil {
			// log.Fatalf("Unable to extract text from page %d: %v", pageNum, err)
			return "", err
		}
		extractedText += text
	}

	// Buat file DOCX baru
	doc := document.New()

	// Tambahkan teks yang diekstraksi ke dalam dokumen DOCX
	doc.AddParagraph().AddRun().AddText(extractedText)

	// Simpan DOCX ke file
	outputPath := namaFile + ".docx"
	err = doc.SaveToFile(outputPath)
	if err != nil {
		// log.Fatalf("Unable to save DOCX file: %v", err)
		return "", err
	}

	// fmt.Println("PDF berhasil dikonversi ke DOCX:", outputPath)
	return "wow berhasil", nil
}
