package webutility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
)

// Dictionary ...
type Dictionary struct {
	my            sync.Mutex
	locales       map[string]map[string]string
	supported     []string
	defaultLocale string
}

// NewDictionary ...
func NewDictionary() *Dictionary {
	return &Dictionary{
		locales: map[string]map[string]string{},
	}
}

// AddTranslations ...
func (d *Dictionary) AddTranslations(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		fName := fileInfo.Name()
		path := directory + "/" + fName
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		loc := stripFileExtension(fName)

		var data interface{}
		err = json.Unmarshal(file, &data)
		if err != nil {
			return err
		}

		l := map[string]string{}
		for k, v := range data.(map[string]interface{}) {
			l[k] = v.(string)
		}

		mu.Lock()
		defer mu.Unlock()
		d.locales[loc] = l
		d.supported = append(d.supported, loc)
	}

	if d.defaultLocale == "" && len(d.supported) > 0 {
		d.defaultLocale = d.supported[0]
	}

	return nil
}

// GetBestMatchLocale ...
func (d *Dictionary) GetBestMatchLocale(req *http.Request) (best string) {
	accepted := d.parseAcceptedLanguages(req.Header.Get("Accept-Language"))

	for i := range accepted {
		if accepted[i].Code == "*" {
			return d.defaultLocale
		}
		for j := range d.supported {
			if accepted[i].Code == d.supported[j] {
				return d.supported[j]
			}
		}
	}

	return d.defaultLocale
}

// Translate ...
func (d *Dictionary) Translate(loc, key string) string {
	return d.locales[loc][key]
}

// SetDefaultLocale ...
func (d *Dictionary) SetDefaultLocale(loc string) error {
	if !d.contains(loc) {
		return fmt.Errorf("locale file not loaded: %s", loc)
	}

	d.defaultLocale = loc

	return nil
}

func (d *Dictionary) contains(loc string) bool {
	for _, v := range d.supported {
		if v == loc {
			return true
		}
	}
	return false
}

// LangWeight ...
type LangWeight struct {
	Code   string
	Weight float64
}

func (d *Dictionary) parseAcceptedLanguages(accepted string) (langs []LangWeight) {
	if accepted == "" {
		langs = append(langs, LangWeight{Code: d.defaultLocale, Weight: 1.0})
		return langs
	}

	var code string
	var weight float64

	parts := strings.Split(accepted, ",")
	for i := range parts {
		parts[i] = strings.Trim(parts[i], " ")

		cw := strings.Split(parts[i], ";")
		paramCount := len(cw)

		if paramCount == 1 {
			// default weight of 1
			code = cw[0]
			weight = 1.0
		} else if paramCount == 2 {
			// parse weight
			code = cw[0]
			weight, _ = strconv.ParseFloat(cw[1][2:], 64)

		}

		langs = append(langs, LangWeight{Code: code, Weight: weight})
	}

	// TODO(marko): sort langs by weights?

	return langs
}

func stripFileExtension(full string) (stripped string) {
	extension := path.Ext(full)
	stripped = strings.TrimSuffix(full, extension)
	return stripped
}
