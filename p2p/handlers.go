package p2p

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func SendFile(rw *bufio.ReadWriter, path string) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// Send filename
	filename := filepath.Base(path)
	_, err = rw.Write([]byte(fmt.Sprintf("%s\n", filename)))
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}

	// Send file size
	fileStats, err := file.Stat()
	if err != nil {
		log.Println(err)
	}
	fileSize := fileStats.Size()
	_, err = rw.Write([]byte(fmt.Sprintf("%d\n", fileSize)))
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}

	// Send file chunks
	hashFunc := sha1.New()
	const chunkSize = 512
	buffer := make([]byte, chunkSize)
	for {
		_, err = file.Read(buffer)
		if err != nil {
			// if err is not EOF, error during reading, else stop
			if err != io.EOF {
				log.Println(err)
			} else {
				break
			}
		}

		_, err = rw.Write(append(buffer, byte('\n')))
		if err != nil {
			log.Println(err)
		}
		err = rw.Flush()
		if err != nil {
			log.Println(err)
		}

		hashFunc.Write(buffer)
	}

	// Send file hash
	hash := hashFunc.Sum(nil)
	_, err = rw.Write(hash)
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}
}

func RecvFile(rw *bufio.ReadWriter) {
	// Recv filename
	filename, err := rw.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	file, err := os.Create(filename)

	// Recv file size

	// Recv file chunks

	// Recv file hash
}
