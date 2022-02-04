package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	file  = getFilePath()
	clips = make(map[string]string)
)

func main() {
	fmt.Printf("Enter master password: ")
	password, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("\nWelcome to govault!")
	fmt.Println("q -> quit, :h -> help")

	if len(password) < 1 || len(password) > 16 {
		fmt.Println("Master password should not be empty and more than 16 characters")
		os.Exit(1)
	}

	password = []byte(fmt.Sprintf("%16v", string(password)))

	read(password)

	for {
		b, _ := terminal.ReadPassword(int(syscall.Stdin))
		s := string(b)

		switch s {
		case ":q":
			os.Exit(0)
		case ":k":
			for k := range clips {
				fmt.Println(k, "****")
			}
			continue
		case ":s":
			for k, v := range clips {
				fmt.Println(k, v)
			}
			continue
		case ":c":
			clearCmd := exec.Command("clear")
			clearCmd.Stdout = os.Stdout
			clearCmd.Run()
			fmt.Println("Welcome to govault!")
		case ":h":
			fmt.Println(":h -> help")
			fmt.Println(":q -> quit")
			fmt.Println(":c -> clear screen")
			fmt.Println(":k -> show keys")
			fmt.Println(":s -> show key/value")
			fmt.Println(":d key -> delete key")
			fmt.Println("key<SPACE>value -> create/update entry")
			fmt.Println("key<ENTER> -> copy value to clipboard")
			continue
		}

		if val, ok := clips[s]; ok {
			copyCmd := exec.Command("pbcopy")
			in, _ := copyCmd.StdinPipe()
			copyCmd.Start()
			in.Write([]byte(val))
			in.Close()
			copyCmd.Wait()
			fmt.Println(s, "copied to clipboard")
		} else {
			idx := strings.Index(s, " ")
			if idx != -1 {
				k := s[0:idx]
				if k == ":d" {
					delete(clips, s[idx+1:])
					fmt.Println(s[idx+1:], "deleted from clips")
				} else {
					clips[k] = s[idx+1:]
					fmt.Println(k, "added to clips")
					save(password)
				}
			}
		}
	}
}

func save(key []byte) {
	var b bytes.Buffer
	writer := io.Writer(&b)

	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(clips); err != nil {
		fmt.Println("Error writing clips to bytes buffer", err)
		return
	}

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher", err)
		return
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(b.Bytes()))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error writing rand bytes to fill iv", err)
		return
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b.Bytes())

	ioutil.WriteFile(file, ciphertext, 0700)
}

func read(key []byte) {
	if _, err := os.Stat(file); err == nil {
		ciphertext, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file")
			return
		}

		// Create the AES cipher
		block, err := aes.NewCipher(key)
		if err != nil {
			fmt.Println("Error creating AES cipher", err)
			return
		}

		// Get the 16 byte IV
		iv := ciphertext[:aes.BlockSize]

		// Remove the IV from the ciphertext
		ciphertext = ciphertext[aes.BlockSize:]

		// Return a decrypted stream
		stream := cipher.NewCFBDecrypter(block, iv)

		// Decrypt bytes from ciphertext
		stream.XORKeyStream(ciphertext, ciphertext)

		decoder := gob.NewDecoder(bytes.NewReader(ciphertext))
		err = decoder.Decode(&clips)

		if err != nil {
			fmt.Println("Error decrypting vault. Please check master password.")
			os.Exit(1)
		}

		fmt.Printf("Vault loaded with %d keys\n", len(clips))
	}
}

func getFilePath() string {
	file := os.Getenv("GO_VAULT_FILE")
	if file == "" {
		file = os.Getenv("HOME") + "/.govault"
	}

	return file
}
