package htpasswd

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

type HashMethod byte

const (
	HashMD5 HashMethod = iota
	HashSha
	HashBcrypt
)

type Entry struct {
	Username     string
	PasswordHash []byte
	HashMethod   HashMethod
}

type File struct {
	users map[string]Entry
}

func (file *File) Users() []string {
	result := make([]string, len(file.users))
	i := 0
	for user := range file.users {
		result[i] = user
		i++
	}
	return result
}

func (file *File) Authenticate(username string, password []byte) (string, bool) {
	if user, ok := file.users[strings.ToLower(username)]; ok {
		if compareHash(password, user.PasswordHash, user.HashMethod) {
			return user.Username, true
		}
	}
	return "", false
}

func LoadHtPasswordFile(filePath string) (*File, error) {
	log.Printf("loading users from [%s]", filePath)
	result := &File{
		users: make(map[string]Entry),
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("the htpassword file was not found. no users will be loaded")
			return result, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		lastColon := bytes.LastIndex(line, []byte(":"))
		user := Entry{
			Username:     string(line[:lastColon]),
			PasswordHash: line[lastColon+1:],
		}
		user.HashMethod = getHashMethod(user.PasswordHash)
		result.users[strings.ToLower(user.Username)] = user
	}
	log.Printf("loaded %d users", len(result.users))
	return result, nil
}
