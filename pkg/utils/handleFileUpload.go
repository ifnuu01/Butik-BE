package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func HandleFileUpload(c echo.Context, fieldName string, destFolder string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", err
	}

	// size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", errors.New("file size exceeds 5MB limit")
	}

	// file extension
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", errors.New("file type not allowed. Allowed: jpg, jpeg, png, gif, webp")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Buat folder jika belum ada
	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return "", err
	}

	// Generate unique filename
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	dstPath := filepath.Join(destFolder, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	urlPath := strings.ReplaceAll(filepath.ToSlash(dstPath), "//", "/")
	baseURL := c.Scheme() + "://" + c.Request().Host
	imageURL := baseURL + "/" + urlPath
	return imageURL, nil
}
