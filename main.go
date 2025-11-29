package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var errEntropyTooLow = errors.New(
	"the entropy of the requested password is too low; please request a password with a larger length",
)

func calculateEntropy(lenCharset, lenPwd int) error {
	minEntropy := 60 // generally acceptable amount of entropy

	// Entropy = log2(lenCharset ** lenPwd)
	entropy := int(math.Round(math.Log2(math.Pow(float64(lenCharset), float64(lenPwd)))))

	if entropy < minEntropy {
		return errEntropyTooLow
	}

	return nil
}

func generatePasswd(lenPwd int) (string, error) {
	chSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcefghijklmopqrstuvwxyz0123456789!@#$%^&*()"
	lenChSet := len(chSet)

	if entropyErr := calculateEntropy(lenChSet, lenPwd); entropyErr != nil {
		return "", entropyErr
	}

	var b strings.Builder
	for i := 0; i <= lenPwd-1; i++ {
		rndIdx, rndErr := rand.Int(rand.Reader, big.NewInt(int64(lenChSet)))
		if rndErr != nil {
			// rand.Int returns an error if rand.Read returns 1. Let's discard the result
			// by decrementing i and trying again.
			i--
			continue
		}

		b.WriteByte(chSet[rndIdx.Int64()])
	}

	return b.String(), nil
}

func printHelp() {
	fmt.Println(`
pwdgen, a minimal password generator

args: desired length of the password
flags:
	- h: help, prints this message
		`)
}

func main() {
	help := flag.Bool("h", false, "prints help text")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments provided.")
		printHelp()
		os.Exit(1)
	}

	if help != nil && *help {
		printHelp()
		os.Exit(0)
	}

	pwdLen, strconvErr := strconv.Atoi(args[0])
	if strconvErr != nil {
		fmt.Println("Invalid argument provided, please provide a number.")
		printHelp()
		os.Exit(1)
	}

	pwd, passwdErr := generatePasswd(pwdLen)
	if passwdErr != nil {
		fmt.Printf("There was an error generating the password: %s\n", passwdErr)
		os.Exit(1)
	}

	if clipboardErr := copyToClipboard(pwd); clipboardErr != nil {
		fmt.Printf("There was an error copying the password to clipboard")
		os.Exit(1)
	}
	os.Exit(0)
}
