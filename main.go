package main
import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var stopAnimation = make(chan bool)
func animate(message string) {
	chars := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-stopAnimation:
			fmt.Printf("\r-> %s   ", message)
			return
		default:
			fmt.Printf("\r-> %s %s", message, chars[i])
			i = (i + 1) % len(chars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
func main() {
	home, _ := os.UserHomeDir()
	mcPath := filepath.Join(home, "AppData", "Roaming", ".minecraft")
	zipURL := "https://github.com/4everdies/bigfun/archive/refs/heads/main.zip"
	tempZip := "skuur.zip"
	fmt.Printf("Path: %s\n", mcPath)
	go animate("Downloading...")
	err := downloadFile(tempZip, zipURL)
	stopAnimation <- true
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		exit()
		return
	}
	fmt.Println("OK!")
	defer os.Remove(tempZip)
	go animate("Loading...")
	foldersToExtract := []string{"mods/", "OneConfig/", "config/"}
	err = unzipSpecific(tempZip, mcPath, foldersToExtract)
	stopAnimation <- true
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		exit()
		return
	}
	fmt.Println("OK!")

	fmt.Println("\nInstalled! :/")
	exit()
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipSpecific(src, dest string, folders []string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		parts := strings.Split(f.Name, "/")
		if len(parts) < 2 {
			continue
		}
		subPath := strings.Join(parts[1:], "/")
		shouldExtract := false
		for _, folder := range folders {
			if strings.HasPrefix(subPath, folder) {
				shouldExtract = true
				break
			}
		}
		if shouldExtract && subPath != "" {
			fpath := filepath.Join(dest, subPath)
			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				return err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func exit() {
	fmt.Println("\nPress Enter to exit")
	fmt.Scanln()
}
