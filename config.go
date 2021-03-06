package dsc

import (
	"strings"

	"github.com/viant/toolbox"
)

//Config represent datastore config.
type Config struct {
	DriverName          string
	PoolSize            int
	MaxPoolSize         int
	Descriptor          string
	Parameters          map[string]string
	SecretParametersURL string //url to JSON object, this delegates credential or secret out of dev
}

//Get returns value for passed in parameter name or panic - please use Config.Has to check if value is present.
func (c *Config) Get(name string) string {
	if result, ok := c.Parameters[name]; ok {
		return result
	}
	panic("Missing value in descriptor " + name)
}

//HasDateLayout returns true if config has date layout, it checks dateLayout or dateFormat parameter names.
func (c *Config) HasDateLayout() bool {
	return toolbox.HasTimeLayout(c.Parameters)
}

//GetDateLayout returns date layout
func (c *Config) GetDateLayout() string {
	return toolbox.GetTimeLayout(c.Parameters)
}

//Has returns true if parameter with passed in name is present, otherwise it returns false.
func (c *Config) Has(name string) bool {
	if _, ok := c.Parameters[name]; ok {
		return true
	}
	return false
}

//Init makes parameter map from encoded parameters if presents, expands descriptor with parameter value using [param_name] matching pattern.
func (c *Config) Init() {
	for key := range c.Parameters {
		macro := "[" + key + "]"
		c.Descriptor = strings.Replace(c.Descriptor, macro, c.Parameters[key], 1)
	}
}

//NewConfig creates new Config, it takes the following parameters
// descriptor - optional datastore connection string with macros that will be looked epxanded from for instance [user]:[password]@[url]
// encodedParameters should be in the following format:   <key1>:<value1>, ...,<keyN>:<valueN>
func NewConfig(driverName string, descriptor string, encodedParameters string) *Config {
	var parameters = toolbox.MakeStringMap(encodedParameters, ":", ",")
	result := &Config{DriverName: driverName, PoolSize: 1, MaxPoolSize: 2, Descriptor: descriptor, Parameters: parameters}
	result.Init()
	return result
}

//NewConfigFromUrl returns new config from url
func NewConfigFromUrl(url string) (*Config, error) {
	resullt := &Config{}
	err := toolbox.LoadConfigFromUrl(url, resullt)
	if err != nil {
		return nil, err
	}
	resullt.Init()
	return resullt, nil
}
