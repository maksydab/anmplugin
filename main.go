package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const port = ":30781"
const zip_name = "plugin.zip"

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	switch os.Args[1] {
	case "serve":
		startserver()
	case "help":
		help()
	default:
		help()
	}
}

func help() {
	fmt.Println("anmplugin cli tool")
	fmt.Println("this tool allows for rapid testing of your plugin without needing to rebundle it every single time")
	fmt.Println()
	fmt.Println("usage for linux:")
	fmt.Println("  ./anmplugin-linux serve   -- start http server on localhost:30781")
	fmt.Println("  ./anmplugin-linux help    -- show this help message")
	fmt.Println("usage for windows:")
	fmt.Println("  ./anmplugin.exe serve     -- start http server on localhost:30781")
	fmt.Println("  ./anmplugin.exe help      -- show this help message")
	fmt.Println("usage for mac:")
	fmt.Println("  ./anmplugin-mac serve     -- start http server on localhost:30781")
	fmt.Println("  ./anmplugin-mac help      -- show this help message")
	fmt.Println()
	fmt.Println("Tip: .ignore behaves like .gitignore")
}

func startserver() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		err := createzip(zip_name)
		if err != nil {
			http.Error(w, "failed to create zip: "+err.Error(), 500)
			return
		}

		http.ServeFile(w, r, zip_name)
	})

	fmt.Println("serving plugin zip at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadignore() map[string]bool {
	ignoremap := make(map[string]bool)

	file, err := os.Open(".ignore")
	if err != nil {
		return ignoremap
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line != "" {
			ignoremap[line] = true
		}
	}

	return ignoremap
}

func isignored(path string, info os.FileInfo, ignore map[string]bool) bool {
	if path == zip_name {
		return true
	}

	base := filepath.Base(path)

	if ignore[path] || ignore[base] {
		return true
	}

	for ignorepath := range ignore {
		// directory ignore support
		if strings.HasSuffix(ignorepath, "/") {
			dir := strings.TrimSuffix(ignorepath, "/")

			if strings.HasPrefix(path, dir+"/") || path == dir {
				return true
			}
		} else {
			if strings.HasPrefix(path, ignorepath+"/") {
				return true
			}

			if info.IsDir() && path == ignorepath {
				return true
			}
		}
	}

	return false
}

func createzip(output string) error {
	ignore := loadignore()

	outfile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outfile.Close()

	zipwriter := zip.NewWriter(outfile)
	defer zipwriter.Close()

	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || isignored(path, info, ignore) {
			return nil
		}

		return addfiletozip(zipwriter, path)
	})
}

func addfiletozip(zipwriter *zip.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = path
	header.Method = zip.Deflate

	writer, err := zipwriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
