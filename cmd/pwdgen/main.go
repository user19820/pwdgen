package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/user19820/pwdgen/internal/args"
	"github.com/user19820/pwdgen/internal/clipboard"
	"github.com/user19820/pwdgen/internal/database"
	"github.com/user19820/pwdgen/internal/encrypt"
	"github.com/user19820/pwdgen/internal/password"
)

//nolint:gocognit // not complicated, just debug logging if statements add more branches
func main() {
	cmd, cmdErr := args.Init()
	if cmdErr != nil {
		if cmd.Debug {
			fmt.Println(cmdErr.Error())
		}

		os.Exit(1)
	}

	switch cmd.Type {
	case args.CmdGen:
		if genPwdErr := genPwd(cmd); genPwdErr != nil {
			if cmd.Debug {
				fmt.Println(genPwdErr.Error())
			}
			os.Exit(1)
		}

		os.Exit(0)
	case args.CmdGet:
		if getPwdErr := getPwd(cmd); getPwdErr != nil {
			if errors.Is(getPwdErr, sql.ErrNoRows) {
				fmt.Printf("no entry found for %s\n", cmd.Name)
				os.Exit(0)
			}

			if cmd.Debug {
				fmt.Println(getPwdErr.Error())
			}
			os.Exit(1)
		}
	case args.CmdInit:
		if initErr := initPwdgen(); initErr != nil {
			if cmd.Debug {
				fmt.Println(initErr.Error())
			}

			os.Exit(1)
		}
	default:
		panic("SHOULD BE UNREACHABLE")
	}
}

func genPwd(args args.Cmd) error {
	pwd, pwdErr := password.Generate(args.Length)
	if pwdErr != nil {
		return pwdErr
	}

	encrypted, encryptErr := encrypt.Encrypt([]byte(pwd))
	if encryptErr != nil {
		return encryptErr
	}

	db, dbHandleErr := database.GetHandle()
	if dbHandleErr != nil {
		return dbHandleErr
	}

	if storeErr := database.Store(db, args.Name, string(encrypted)); storeErr != nil {
		return storeErr
	}

	return clipboard.Copy(pwd)
}

func getPwd(args args.Cmd) error {
	db, dbHandleErr := database.GetHandle()
	if dbHandleErr != nil {
		return dbHandleErr
	}

	pwd, pwdErr := database.Retrieve(db, args.Name)
	if pwdErr != nil {
		return pwdErr
	}

	decrypted, decryptErr := encrypt.Decrypt([]byte(pwd))
	if decryptErr != nil {
		return decryptErr
	}

	return clipboard.Copy(string(decrypted))
}

func initPwdgen() error {
	if err := database.Init(); err != nil {
		return err
	}

	return encrypt.Setup()
}
