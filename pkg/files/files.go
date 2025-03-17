package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const initialFileContent = `
# %s

## %s


## %s


## %s


`

// checks if the file is present in the current dir.
func CheckFile(fileName string) bool {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filePath := dir + fileName
	_, err = os.Stat(filePath)
	return err == nil
}

// create a file
func CreateFile(fileName string) {
	f, err := OpenFileWriteOnly(fileName)
	if err != nil {
		log.Fatalf("Cannot create file %q, Err: %v", fileName, err)
	}
	defer f.Close()
}

// writes the initial content to the file
func WriteInitialContent(fileName string) {
	f, err := OpenFileWriteOnly(fileName)
	if err != nil {
		log.Fatalf("Cannot create file %q, Err: %v", fileName, err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(initialFileContent, "", "TODO", "DOING", "DONE"))
	if err != nil {
		log.Fatalf("Cannot write to the file (%v): %v", fileName, err)
	}
}

// opens a file in write only mode
func OpenFileWriteOnly(fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", err)
	}
	return os.OpenFile(dir+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// returns the name of the current working directory
func GetDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get directory name: %v", err)
	}
	dirs := strings.Split(dir, "/")
	dirName := dirs[len(dirs)-1]
	return dirName
}

// returns the complete path of the file, given the filename
func FilePath(fileName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filePath := dir + fileName
	return filePath, nil
}

// opens the given file and returns the content of the file
// line by line as a slice of string
func OpenFile(fileName string) []string {
	filePath, err := FilePath(fileName)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileContent []string
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}
	return fileContent
}

// writes a slice of string to a file line by line
func WriteFile(fileContent []string, fileName string) error {
	filePath, err := FilePath(fileName)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range fileContent {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
