package controllers

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func unzipWithPassword(src, password, dest string) error {
	// Buka file ZIP
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Buat folder tujuan jika belum ada
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	// Iterasi melalui semua file dalam ZIP
	for _, f := range reader.File {
		// Atur password untuk file yang dienkripsi
		f.SetPassword(password)

		filePath := filepath.Join(dest, f.Name)

		// Cek apakah file adalah direktori
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// Jika bukan direktori, buat file
		err = extractFileWithPassword(f, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFileWithPassword(f *zip.File, filePath string) error {
	// Buka file dalam ZIP
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Buat direktori jika file berada di dalam folder
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// Buat file baru di sistem
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Salin isi file dari ZIP ke file baru
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	return nil
}