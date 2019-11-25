package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
)

type Configuration struct {
	Key          string
	Token        string
	AccessToken  string
	AccessSecret string
}

func ReadConfiguration(configurationFile string) *Configuration {

	if len(configurationFile) == 0 {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		configurationFile = usr.HomeDir + "/.twitter_cli"
	}

	data, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		panic(err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(data, &dat); err != nil {
		return &Configuration{}
	}
	result := Configuration{
		dat["key"].(string),
		dat["token"].(string),
		dat["accessToken"].(string),
		dat["accessSecret"].(string)}

	return &result
}

func (c *Configuration) SaveConfiguration() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	dat := map[string]string{"key": c.Key, "token": c.Token, "accessToken": c.AccessToken, "accessSecret": c.AccessSecret}

	bytes, _ := json.Marshal(dat)
	e := ioutil.WriteFile(usr.HomeDir + "/.twitter_cli", bytes, 0600)
	if e != nil {
		panic(e)
	}
}