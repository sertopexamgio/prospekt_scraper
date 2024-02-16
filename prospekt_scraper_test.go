package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type mockReaderCloser struct {
	io.Reader
}

func (mockReaderCloser) Close() error { return nil }

func TestSavePDF(t *testing.T) {
	pdfContent := []byte("Sample PDF content")
	reader := bytes.NewReader(pdfContent)

	savePDF(mockReaderCloser{reader}, "test.pdf")

	savedContent, err := os.ReadFile("test.pdf")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(savedContent, pdfContent) {
		t.Error("Saved PDF content does not match original PDF content")
	}
}
