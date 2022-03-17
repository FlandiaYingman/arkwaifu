package res

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	KeySize         = 16
	IVSize          = 16
	MagicOffsetSize = 128
)

var (
	chatMask = make([]byte, 32)
)

func SetChatMask(cm []byte) {
	copy(chatMask, cm)
}

func decrypt(ctx context.Context, srcDir string, dstDir string) error {
	return filepath.WalkDir(srcDir, func(srcPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(srcPath, ".bytes") {
			return nil
		}

		err = decryptFile(srcPath, srcDir, dstDir)
		if err != nil {
			logrus.Error(err)
			return nil
		}

		return nil
	})
}

func decryptFile(srcPath string, srcDir string, dstDir string) error {
	hasMagicOffset, _ := filepath.Match("**/levels/**", filepath.ToSlash(srcPath))
	hasMagicOffset = !hasMagicOffset

	if filepath.Base(srcPath) == "story_review_table.bytes" {
		print("!!")
	}
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		return errors.WithStack(err)
	}
	cm := make([]byte, len(chatMask))
	copy(cm, chatMask)
	plainText, err := decryptData(srcContent, cm, hasMagicOffset)
	if err != nil {
		return errors.WithStack(err)
	}

	switch {
	case isBson(plainText):
		bsonContent := plainText
		err := saveFile(srcPath, srcDir, dstDir, ".bson", bsonContent)
		if err != nil {
			return err
		}

		var bsonObj bson.M
		err = bson.Unmarshal(bsonContent, &bsonObj)
		if err != nil {
			return errors.WithStack(err)
		}
		jsonContent, err := json.Marshal(bsonObj)
		err = saveFile(srcPath, srcDir, dstDir, ".json", jsonContent)
		if err != nil {
			return err
		}
	case isJson(plainText):
		err := saveFile(srcPath, srcDir, dstDir, ".json", plainText)
		if err != nil {
			return err
		}
	case strings.HasSuffix(srcPath, ".lua.bytes"):
		err := saveFile(srcPath, srcDir, dstDir, "", plainText)
		if err != nil {
			return err
		}
	default:
		err := saveFile(srcPath, srcDir, dstDir, ".txt", plainText)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveFile(srcPath string, srcDir string, dstDir string, dstExt string, content []byte) error {
	dstPath := pathutil.ReplaceParentExt(srcPath, srcDir, dstDir, dstExt)

	var err error
	err = os.MkdirAll(filepath.Dir(dstPath), 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.WriteFile(dstPath, content, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func isBson(bytes []byte) bool {
	var m bson.M
	err := bson.Unmarshal(bytes, &m)
	return err == nil
}
func isJson(bytes []byte) bool {
	var m map[string]interface{}
	err := json.Unmarshal(bytes, &m)
	return err == nil
}

// decryptData decrypts the
func decryptData(cipherText []byte, chatMask []byte, hasMagicOffset bool) (result []byte, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if recErr, ok := rec.(error); ok {
				err = recErr
			} else {
				err = errors.Errorf("%+v", recErr)
			}
		}
	}()

	if hasMagicOffset {
		cipherText = cipherText[MagicOffsetSize:]
	}
	key := extractKey(chatMask)
	iv := extractIV(cipherText, chatMask)
	cipherText = cipherText[IVSize:]
	cipherText = decryptAesCbc(cipherText, key, iv)
	cipherText = unpadAnsiX923(cipherText)
	return cipherText, nil
}

func extractKey(chatMask []byte) []byte {
	return chatMask[:KeySize]
}
func extractIV(raw []byte, chatMask []byte) []byte {
	iv := chatMask[len(chatMask)-IVSize:]
	for i := range iv {
		iv[i] ^= raw[i]
	}
	return iv
}

// decryptAesCbc decrypts cipher text with AES-CBC method.
func decryptAesCbc(cipherText []byte, key []byte, iv []byte) []byte {
	// Check cipherText size is whether a multiple of AES block size
	cipherTextSize := len(cipherText)
	aesBlockSize := aes.BlockSize
	if cipherTextSize%aesBlockSize != 0 {
		panic(errors.Errorf("cipherText size %v is not a multiple of AES block size %v", cipherTextSize, aesBlockSize))
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cbcCipher := cipher.NewCBCDecrypter(aesCipher, iv)
	cbcCipher.CryptBlocks(cipherText, cipherText)

	// Now cipherText become plainText
	plainText := cipherText
	return plainText
}

// unpadAnsiX923 unpads a decrypted plain text with ANSI X9.23 method.
func unpadAnsiX923(plainText []byte) []byte {
	plainTextSize := len(plainText)
	paddingSize := int(plainText[plainTextSize-1])
	plainText = plainText[:plainTextSize-paddingSize]
	return plainText
}
