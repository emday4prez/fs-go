package file

import (
	"errors"
	"io"
	"mime/multipart"
	"reflect"
	"testing"
)

// MockStorage is our fake implementation of the Storage interface for testing.
// It lives in the same package, so it can be used by the test.
type MockStorage struct {
	SaveFunc   func(fileHeader *multipart.FileHeader) (string, error)
	ListResult []string
	ListError  error
	GetResult  io.ReadCloser
	GetPath    string
	GetError   error
}

// Implement the Storage interface for our MockStorage.
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

// TestFileService_ListFiles is our first unit test.
func TestFileService_ListFiles(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{
		ListResult: []string{"test1.txt", "image.png"},
		ListError:  nil,
	}
	fileService := NewFileService(mockStorage)

	// Act
	actualFilenames, err := fileService.ListFiles()

	// Assert
	if err != nil {
		t.Errorf("ListFiles() returned an unexpected error: %v", err)
	}
	expectedFilenames := []string{"test1.txt", "image.png"}
	if !reflect.DeepEqual(actualFilenames, expectedFilenames) {
		t.Errorf("ListFiles() returned %v, want %v", actualFilenames, expectedFilenames)
	}
}

// TestFileService_ListFiles_Error tests the error case.
func TestFileService_ListFiles_Error(t *testing.T) {
	// Arrange
	mockStorage := &MockStorage{
		ListResult: nil,
		ListError:  errors.New("storage is down"),
	}
	fileService := NewFileService(mockStorage)

	// Act
	_, err := fileService.ListFiles()

	// Assert
	if err == nil {
		t.Errorf("ListFiles() expected an error, but got nil")
	}
}

// TestFileService_SaveFile verifies that the SaveFile service method
// correctly calls the underlying storage's Save method.
func TestFileService_SaveFile(t *testing.T) {
	// Arrange
	var saveCalled bool
	mockStorage := &MockStorage{
		SaveFunc: func(fileHeader *multipart.FileHeader) (string, error) {
			saveCalled = true
			return "mock/path", nil
		},
	}
	fileService := NewFileService(mockStorage)
	dummyFileHeader := &multipart.FileHeader{}

	// Act
	err := fileService.SaveFile(dummyFileHeader)

	// Assert
	if err != nil {
		t.Errorf("SaveFile() returned an unexpected error: %v", err)
	}
	if !saveCalled {
		t.Errorf("Expected storage.Save() to be called, but it was not.")
	}
}
