package bart

import (
	"fmt"
	"io/ioutil"
)

func ProcessFile(templateFilename string, send bool, c *Config) error {
	// FIXME: Disentangle this mess

	data, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return err
	}
	data_s := string(data)

	if send {
		fmt.Printf("Please enter your credentials for \"%s\"\n", c.EmailServer.Hostname)
		ap := new(authPair)
		ap.prompt()

		for recipient, context := range c.Recipients {
			for k, v := range c.GlobalContext {
				if _, present := context[k]; !present {
					context[k] = v
				}
			}

			email := NewEmail().AddAuthor(&c.Author).AddRecipient(recipient).AddContent(data_s).Build(context)

			fmt.Printf("Will send to %v\n", email.GetRecipients())
			if err := email.Send(&c.EmailServer, ap); err != nil {
				return err
			}
		}
	} else {
		for recipient, context := range c.Recipients {
			for k, v := range c.GlobalContext {
				if _, present := context[k]; !present {
					context[k] = v
				}
			}

			email := NewEmail().AddAuthor(&c.Author).AddRecipient(recipient).AddContent(data_s).Build(context)

			fmt.Printf("Send flag not set: opening preview in \"%s\"\n", c.Author.Browser)
			if err := email.OpenInBrowser(c.Author.Browser); err != nil {
				return err
			}
		}
	}
	return nil
}
