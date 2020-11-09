package key

import (
	"fmt"
	"io/ioutil"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
)

//Verify the signature of the signed message
func Verify() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		publicKey := c.String(cli.StringArg{
			Name: "PUBLICKEY",
			Desc: "Public key",
		})
		signature := c.String(cli.StringArg{
			Name: "SIGNATURE",
			Desc: "Signature of the message",
		})
		message := c.String(cli.StringOpt{
			Name: "m message",
			Desc: "Message to be verified",
		})
		messageFile := c.String(cli.StringOpt{
			Name: "f messagefile",
			Desc: "Message file to be verified",
		})

		c.Spec = "PUBLICKEY SIGNATURE [-m=<message to be verified>] | [-f=<Message File to be verified>]"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			var msg []byte
			var err error
			if *message != "" {
				msg = []byte(*message)
			} else if *messageFile != "" {
				msg, err = ioutil.ReadFile(*messageFile)
				if err != nil {
					cmd.PrintErrorMsg("Failed to read the file: %v", err)
					return
				}
			} else {
				cmd.PrintWarnMsg("Please enter a message to verify.")
				c.PrintHelp()
				return
			}

			var sign crypto.Signature
			publickey, err := crypto.PublicKeyFromString(*publicKey)
			if err != nil {
				cmd.PrintErrorMsg("%v", err)
				return
			}
			sign, err = crypto.SignatureFromString(*signature)
			if err != nil {
				cmd.PrintErrorMsg("%v", err)
				return
			}

			fmt.Println()
			verify := publickey.Verify(msg, &sign)
			if verify {
				cmd.PrintSuccessMsg("Signature is verified successfully!")
			} else {
				cmd.PrintErrorMsg("Signature verification failed!")
			}
		}
	}
}
