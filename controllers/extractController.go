package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Handler untuk menerima file ZIP
func UploadAndExtractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Hanya mendukung metode POST", http.StatusMethodNotAllowed)
		return
	}

	// Parsing form untuk mendapatkan file
	err := r.ParseMultipartForm(10 << 20) // Batas 10MB
	if err != nil {
		http.Error(w, "Gagal memproses file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Ambil file dari form-data
	file, fileHeader, err := r.FormFile("zipfile")
	if err != nil {
		http.Error(w, "Gagal membaca file dari request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file ZIP sementara di server
	tempZipPath := "./temp/" + fileHeader.Filename
	os.MkdirAll("./temp", os.ModePerm)
	outFile, err := os.Create(tempZipPath)
	if err != nil {
		http.Error(w, "Gagal membuat file sementara", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Gagal menyimpan file sementara", http.StatusInternalServerError)
		return
	}

	// Direktori tujuan ekstraksi
	extractDir := "./output/" + fileHeader.Filename
	os.MkdirAll(extractDir, os.ModePerm)

	// Ekstrak file ZIP
	err = extractZip(tempZipPath, extractDir)
	if err != nil {
		http.Error(w, "Gagal mengekstrak file ZIP", http.StatusInternalServerError)
		return
	}

	// Hapus file ZIP sementara
	os.Remove(tempZipPath)

	// Balas ke client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File ZIP berhasil diekstrak ke: %s", extractDir)
}

// Fungsi untuk ekstraksi file ZIP
func extractZip(zipPath, destination string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(destination, f.Name)

		// Buat direktori jika diperlukan
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// Ekstrak file
		if err := extractFile(f, filePath); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, filePath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Buat file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}
