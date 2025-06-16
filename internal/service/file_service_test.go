package service

import (
	"errors"
	"io"
	"mime/multipart"
	"reflect"
	"testing"
)

type MockStorage struct {
	SaveFunc   func(fileHeader *multipart.FileHeader) (string, error)
	ListResult []string
	ListError  error
	GetResult  io.ReadCloser
	GetPath    string
	GetError   error
}

func (m *MockStorage) Save(fileHeader *multipart.FileHeader) (string, error) {

	if m.SaveFunc != nil {
		return m.SaveFunc(fileHeader)
	}
	return "mock_path/mock_file.txt", nil
}

func (m *MockStorage) List() ([]string, error) {

	return m.ListResult, m.ListError
}

func (m *MockStorage) Get(filename string) (io.ReadCloser, string, error) {
	return m.GetResult, m.GetPath, m.GetError
}

func TestFileService_ListFiles(t *testing.T) {

	mockStorage := &MockStorage{

		ListResult: []string{"test1.txt", "image.png"},
		ListError:  nil,
	}

	fileService := NewFileService(mockStorage)

	actualFilenames, err := fileService.ListFiles()

	if err != nil {
		t.Errorf("ListFiles() returned an unexpected error: %v", err)
	}

	expectedFilenames := []string{"test1.txt", "image.png"}
	if !reflect.DeepEqual(actualFilenames, expectedFilenames) {
		t.Errorf("ListFiles() returned %v, want %v", actualFilenames, expectedFilenames)
	}
}

func TestFileService_ListFiles_Error(t *testing.T) {
	mockStorage := &MockStorage{

		ListError: errors.New("storage is down"),
	}
	fileService := NewFileService(mockStorage)

	_, err := fileService.ListFiles()

	if err == nil {
		t.Errorf("ListFiles() expected an error, but got nil")
	}
}
