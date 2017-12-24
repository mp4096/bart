package bart

import (
	"fmt"
	"io/ioutil"
)

/// Merge two maps. Key-value pairs from `secondary` are added to primary.
/// If a key is present in both maps, the `primary`'s value is retained.
func mergeMaps(primary map[string]string, secondary map[string]string) {
	for k, v := range secondary {
		if _, present := primary[k]; !present {
			primary[k] = v
		}
	}
}

func ProcessFile(templateFilename string, send bool, c *Config) error {
	data, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return err
	}
	data_s := string(data)

	ap := new(authPair)
	if send {
		fmt.Printf("Please enter your credentials for \"%s\"\n", c.EmailServer.Hostname)
		ap.prompt()
	}

	for recipient, localContext := range c.Recipients {
		// Add global context to local context
		mergeMaps(localContext, c.GlobalContext)

		email := NewEmail().AddAuthor(&c.Author).AddRecipient(recipient).AddContent(data_s).Build(localContext)

		if send {
			fmt.Printf("Will send to %v\n", email.GetRecipients())
			if err := email.Send(&c.EmailServer, ap); err != nil {
				return err
			}
		} else {
			fmt.Printf("Send flag not set: opening preview in \"%s\"\n", c.Author.Browser)
			if err := email.OpenInBrowser(c.Author.Browser); err != nil {
				return err
			}
		}
	}

	return nil
}
