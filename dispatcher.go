package bart

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
)

// Matches an underscore at the beginning of the string
var underscoreRegex = regexp.MustCompile(`^\_`)

// mergeMaps merges two maps. Key-value pairs from `secondary` are added to primary.
// If a key is present in both maps, the `primary`'s value is retained.
func mergeMaps(primary map[string]string, secondary map[string]string) {
	for k, v := range secondary {
		if _, present := primary[k]; !present {
			primary[k] = v
		}
	}
}

/// Verify a context. Key 'subject' must be present.
func checkForSubject(context map[string]string) error {
	if _, ok := context["subject"]; ok {
		return nil
	}
	return errors.New("no 'subject' tag found. It must be contained in the context")
}

/// Verify a context. No keys beginning with an underscore are allowed.
func checkForNoUnderscores(context map[string]string) error {
	for k := range context {
		if underscoreRegex.MatchString(k) {
			return errors.New("tags beginning with an underscore are not allowed")
		}
	}
	return nil
}

func ProcessFile(templateFilename string, send bool, c *Config) error {
	data, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return err
	}
	dataAsString := string(data)

	ap := new(authPair)
	if send {
		fmt.Printf("Please enter your credentials for \"%s\"\n", c.EmailServer.Hostname)
		ap.prompt()
	}

	for recipient, localContext := range c.Recipients {
		// Add global context to local context
		mergeMaps(localContext, c.GlobalContext)

		if err := checkForSubject(localContext); err != nil {
			return err
		}
		if err := checkForNoUnderscores(localContext); err != nil {
			return err
		}

		localContext["__subject_encoded__"] = EncodeRfc1342(localContext["subject"])

		email := NewEmail().AddAuthor(&c.Author).AddRecipient(recipient).AddContent(dataAsString).Build(localContext)

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
