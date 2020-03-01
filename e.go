package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Profile struct {
	Name   string
	Values ProfileValues
}

type ProfileValues map[string]string

type E struct {
	Selected string
	Profiles map[string]*Profile
}

var cfgdir string
var pointerPath string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	cfgdir = path.Join(homedir, ".e")
	pointerPath = path.Join(cfgdir, ".profile")
}

func createPointerFile() error {
	err := ioutil.WriteFile(pointerPath, []byte(""), 0777)
	if err != nil {
		return err
	}
	return nil
}

func bootstrap() error {
	if _, err := os.Stat(pointerPath); err != nil {
		err := createPointerFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) GetValues() (ProfileValues, error) {
	if p.Values == nil {
		err := p.ReadValues()
		if err != nil {
			return nil, err
		}
	}
	return p.Values, nil
}

func (p *Profile) ReadValues() error {
	b, err := ioutil.ReadFile(path.Join(cfgdir, p.Name))
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	p.Values = make(ProfileValues)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		i := strings.Index(line, "=")
		if i == -1 {
			continue
		}
		key := line[:i]
		value := line[i+1:]
		p.Values[key] = value
	}
	return nil
}

func (p *Profile) SetValues() error {
	values, err := p.GetValues()
	if err != nil {
		return err
	}
	content := ""
	for k, v := range values {
		if v == "" {
			continue
		}
		content += fmt.Sprintf("%s=%s\n", k, v)
	}
	err = ioutil.WriteFile(path.Join(cfgdir, p.Name), []byte(content), 0777)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) SetValue(key string, value string) error {
	values, err := p.GetValues()
	if err != nil {
		return err
	}
	values[key] = value
	err = p.SetValues()
	if err != nil {
		return err
	}

	return nil
}

func (e *E) CreateProfile(name string) (*Profile, error) {
	profilePath := path.Join(cfgdir, name)
	if _, err := os.Stat(profilePath); err == nil {
		return nil, fmt.Errorf("Could not create profile: Profile %s already exists", name)
	}
	err := ioutil.WriteFile(profilePath, []byte(""), 0777)
	if err != nil {
		return nil, err
	}
	profile := &Profile{Name: name}
	e.Profiles[name] = profile
	return profile, nil
}

func NewE() *E {
	return &E{}
}

func ProfileFromFile(file os.FileInfo) (*Profile, error) {
	profile := &Profile{Name: file.Name()}
	return profile, nil
}

func (e *E) SelectProfile(name string) error {
	_, ok := e.Profiles[name]
	if !ok {
		return fmt.Errorf("Could not select profile: Profile %s does not exists", name)
	}
	e.Selected = name
	err := ioutil.WriteFile(path.Join(cfgdir, ".profile"), []byte(name), 0777)
	if err != nil {
		return err
	}
	return nil
}

func Load() (*E, error) {
	err := os.MkdirAll(cfgdir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = bootstrap()
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(cfgdir)
	if err != nil {
		return nil, err
	}
	profiles := make(map[string]*Profile)
	for _, f := range files {
		if f.Name() == ".profile" {
			continue
		}
		profile, err := ProfileFromFile(f)
		if err != nil {
			return nil, err
		}
		profiles[profile.Name] = profile
	}
	e := NewE()
	e.Profiles = profiles
	b, err := ioutil.ReadFile(path.Join(cfgdir, ".profile"))
	if err != nil {
		return nil, err
	}
	e.Selected = string(b)
	return e, nil
}

func (e *E) GetProfile() *Profile {
	return e.Profiles[e.Selected]
}
