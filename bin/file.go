package bin

import (
	b64 "encoding/base64"
	"github.com/Scarsz/bincli/crypto"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	UUID        uuid.UUID
	Name        string
	Content     []byte
	Description string
}

func FileFromFileName(fileName string) File {
	file, err := os.Open(fileName)
	if err != nil { panic(err) }

	data, err := ioutil.ReadAll(file)
	if err != nil { panic(err) }

	return File {
		UUID:    uuid.New(),
		Name:    fileName,
		Content: data,
	}
}

func FileFromText(name string, text string, description string) File {
	return File{
		UUID:        uuid.New(),
		Name:        name,
		Content:     []byte(text),
		Description: text,
	}
}

func (file *File) ContentType() string {
	ext := strings.ToLower(filepath.Ext(file.Name))

	switch ext {
	case ".log":
	case ".txt":
	case ".rtf":
		return "text/plain"
	}

	return strings.Split(http.DetectContentType(file.Content[0:512]), ";")[0]
}

func (file *File) ContentText() string {
	return string(file.Content)
}

func (file *File) EncryptAndEncode(key []byte) (name, content, contentType, description string) {
	name = b64.StdEncoding.EncodeToString(crypto.EncryptString(key, file.Name))
	content = b64.StdEncoding.EncodeToString(crypto.Encrypt(key, file.Content))
	contentType = b64.StdEncoding.EncodeToString(crypto.EncryptString(key, file.ContentType()))
	description = b64.StdEncoding.EncodeToString(crypto.EncryptString(key, file.Description))
	return name, content, contentType, description
}

func (file *File) SerializeMap(key []byte) map[string]interface{} {
	m := make(map[string]interface{})

	name, content, contentType, description := file.EncryptAndEncode(key)
	m["name"] = name
	m["content"] = content
	m["type"] = contentType
	if file.Description != "" {
		m["description"] = description
	}

	return m
}

