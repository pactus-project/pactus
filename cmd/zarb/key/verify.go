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
		publicKeyArg := c.String(cli.StringArg{
			Name: "PUBLICKEY",
			Desc: "Public key",
		})
		signatureArg := c.String(cli.StringArg{
			Name: "SIGNATURE",
			Desc: "Signature of the message",
		})
		messageOpt := c.String(cli.StringOpt{
			Name: "m message",
			Desc: "Message to verify",
		})
		messageFileOpt := c.String(cli.StringOpt{
			Name: "f messagefile",
			Desc: "Message file to verify",
		})

		c.Spec = "PUBLICKEY SIGNATURE"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			var msg []byte
			var err error
			if *messageOpt != "" {
				msg = []byte(*messageOpt)
			} else if *messageFileOpt != "" {
				msg, err = ioutil.ReadFile(*messageFileOpt)
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
			publickey, err := crypto.PublicKeyFromString(*publicKeyArg)
			if err != nil {
				cmd.PrintErrorMsg("%v", err)
				return
			}
			sign, err = crypto.SignatureFromString(*signatureArg)
			if err != nil {
				cmd.PrintErrorMsg("%v", err)
				return
			}

			cmd.PrintLine()
			verify := publickey.Verify(msg, sign)
			if verify {
				cmd.PrintSuccessMsg("Signature is verified successfully!")
			} else {
				cmd.PrintErrorMsg("Signature verification failed!")
			}
		}
	}
}
