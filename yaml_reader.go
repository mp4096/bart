package bart

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Author struct {
	Name    string
	Email   string
	Browser string
}

type EmailServer struct {
	Hostname string
	Port     int
}

// Recipients contains emails and their respective contexts.
// First map: Emails as keys
// Second map: Mustache context
type Recipients map[string]map[string]string

type Config struct {
	EmailServer   EmailServer       `yaml:"email_server"`
	GlobalContext map[string]string `yaml:"global_context"`
	Author        Author
	Recipients    Recipients
}

func (t *Config) ImportFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.UnmarshalStrict(data, &t)
}
