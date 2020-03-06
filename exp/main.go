package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/sliceking/galleria/hash"
)

func main() {
	toHash := []byte("this string to hash")
	h := hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(base64.URLEncoding.EncodeToString(b))

	hmac := hash.NewHMAC("my-secret-key")
	fmt.Println(hmac.Hash("this string to hash"))

}
