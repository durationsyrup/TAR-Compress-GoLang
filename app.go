package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func logFatal(e error) {
	// Just because
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	compress()
}

func compress() {
	// Configure this how ever you want
	var c, o string
	fmt.Print("\nWhat would you like to compress? ")
	fmt.Scanf("%s", &c)
	fmt.Print("\nWhat would you like to name the output file? ")
	fmt.Scanf("%s", &o)

	files, _ := readAllDir(c)
	compressFiles(files, o)
}

func compressFiles(list []fileInfo, output string) {
	var totalSize int64 = 0
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, f := range list {
		err := tw.WriteHeader(&tar.Header{
			Name: f.Path,
			Mode: 0600,
			Size: f.File.Size(),
		})
		logFatal(err)
		data, err := ioutil.ReadFile(f.Path)
		logFatal(err)
		_, err = tw.Write(data)
		logFatal(err)
		totalSize += f.File.Size()
	}
	tw.Close()

	if len(output) >= 4 {
		if output[len(output)-4:] == ".tar" {
			output = output[:len(output)-4]
		}
	}

	ioutil.WriteFile(output+".tar", buf.Bytes(), 0644)

	f, err := os.Open(output + ".tar")
	defer f.Close()
	logFatal(err)
	info, err := f.Stat()
	logFatal(err)
	fmt.Printf("\nFinished! %s.tar was created, Original Size: %v bytes / Compressed Size: %v bytes.\n", output, totalSize, info.Size())
}

// returns List of File Info, List of Folder Info
func readAllDir(filename string) ([]fileInfo, []fileInfo) {
	// There are multiples ways to do this
	folder, err := os.Open(filename)
	logFatal(err)
	defer folder.Close()
	if info, err := folder.Stat(); !info.IsDir() {
		logFatal(err)
		return []fileInfo{
			fileInfo{Path: filename, File: info},
		}, []fileInfo{}
	}

	files, err := folder.Readdir(0)
	logFatal(err)

	Directories := []fileInfo{}
	Files := []fileInfo{}

	for _, f := range files {
		if f.IsDir() {
			Directories = append(Directories, fileInfo{
				Path: filename + "/" + f.Name(),
				File: f,
			})
		} else {
			Files = append(Files, fileInfo{
				Path: filename + "/" + f.Name(),
				File: f,
			})
		}
	}

	for i := 0; i < len(Directories); i++ {
		if Directories[i].File.IsDir() {
			folder, err := os.Open(Directories[i].Path)
			logFatal(err)
			defer folder.Close()
			files, err := folder.Readdir(0)
			logFatal(err)
			for _, f := range files {
				if f.IsDir() {
					Directories = append(Directories, fileInfo{
						Path: Directories[i].Path + "/" + f.Name(),
						File: f,
					})
				} else {
					Files = append(Files, fileInfo{
						Path: Directories[i].Path + "/" + f.Name(),
						File: f,
					})
				}
			}
		}
	}
	return Files, Directories
}

type fileInfo struct {
	Path string
	File os.FileInfo
}
