package common

import (
	"errors"
	"github.com/hhy5861/go-common/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Data interface{}
}

func NewConfig(data interface{}) *Config {
	return &Config{Data: data}
}

func (c *Config) Read(file string) {
	var err error

	if err = c.ReadYaml(file); err != nil {
		logger.Error(err)
	}

	if err = c.ReadEnv(); err != nil {
		logger.Error(err)
	}
}

func (c *Config) ReadYaml(file string) error {
	var (
		err  error
		data []byte
	)

	data, err = ioutil.ReadFile(file)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = yaml.Unmarshal([]byte(data), c.Data)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Config) ReadEnv() error {
	var (
		elem reflect.Value
		err  error
	)

	elem = reflect.ValueOf(c.Data).Elem()
	err = reflectEnvAction("", elem)
	if err != nil {
		return err
	}

	return nil
}

func reflectEnvAction(parent string, value reflect.Value) error {
	var (
		err   error
		count int
		tag   string
		to    reflect.Type
	)

	to = value.Type()
	count = value.NumField()

	for i := 0; i < count; i++ {
		f := value.Field(i)

		if f.IsValid() && f.CanSet() {
			tag = string(to.Field(i).Tag.Get("yaml"))
			envStr := ""
			if tag != "" {
				envStr = parent + strings.ToUpper(tag)
			} else {
				envStr = parent + strings.ToUpper(to.Field(i).Name)
			}

			env := os.Getenv(envStr)

			if env != "" {
				switch f.Kind() {
				case reflect.String:
					f.SetString(env)
				case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
					d, err := strconv.ParseInt(env, 10, 64)
					if err != nil {
						logger.Error(err)
						return err
					}
					switch f.Kind() {
					case reflect.Int:
						f.Set(reflect.ValueOf(int(d)))
					case reflect.Uint:
						f.Set(reflect.ValueOf(uint(d)))
					case reflect.Int64:
						f.Set(reflect.ValueOf(int64(d)))
					case reflect.Uint64:
						f.Set(reflect.ValueOf(uint64(d)))
					}
				case reflect.Float64:
					d, err := strconv.ParseFloat(env, 64)
					if err != nil {
						logger.Error(err)
						return err
					}
					f.SetFloat(d)
				case reflect.Bool:
					if strings.ToLower(env) == "true" {
						f.SetBool(true)
					} else {
						f.SetBool(false)
					}
				}
			}

			switch f.Kind() {
			case reflect.Struct:
				err = reflectEnvAction(envStr+"_", f)
				if err != nil {
					return err
				}
			case reflect.Ptr:
				err = reflectEnvAction(envStr+"_", f)
				if err != nil {
					return err
				}
			}

			if err != nil {
				return errors.New("Parsing" + tag + "failed")
			}
		}
	}

	return nil
}
