package loadkit

import (
	"fmt"
	"log/slog"
	"os"
	p "path"
	"strings"

	"github.com/goccy/go-yaml"
	m "github.com/raspidrum-srv/internal/model"
	f "github.com/raspidrum-srv/internal/repo/file"
)

/* Временное решение по преобразованию массивов controls и layers в map
Исходный вид:
```yaml
  layers:
    - name: Closed
      midiKey: HHCLOSEKEY
      controls:
        - name: Volume
          type: volume
          key: HATCV
    - name: Open
      midiKey: HHOPENKEY
      controls:
        - name: Volume
          type: volume
          key: HATOV
		- name: Foot open
      midiKey: HHFOOTOPKEY
      controls:
        - name: Volume
          type: volume
          key: HATFOV
```
 Требуемое:
 ```yaml
  layers:
    closed:
			midiKey: HHCLOSEKEY
      controls:
        volume:
					name: Volume #optional
          type: volume
          key: HATCV
    open:
      midiKey: HHOPENKEY
      controls:
        volume:
          type: volume
          key: HATOV
		foot_open:
      name: Foot open #optional
      midiKey: HHFOOTOPKEY
      controls:
        volume:
          type: volume
          key: HATFOV
```
TODO:
 - преобразование (Layer|Control).name в key:
		- приведение к нижнему регистру
		- замена " " на "_"
 - если name содержит пробел, то добавлять опциональный атрибут name
 - изменить схему
	- instrument Control, Layer
	- kit_preset Control, Layer
 - перезалить
	- kit
	- kit_preset
*/

type NewInstrument struct {
	Id          int64                `yaml:"-"`
	Uid         string               `yaml:"UUID"`
	Key         string               `yaml:"key"`
	Name        string               `yaml:"name"`
	FullName    string               `yaml:"fullName,omitempty"`
	Type        string               `yaml:"type"`
	SubType     string               `yaml:"subtype"`
	Description string               `yaml:"description,omitempty"`
	Copyright   string               `yaml:"copyright,omitempty"`
	Licence     string               `yaml:"licence,omitempty"`
	Credits     string               `yaml:"credits,omitempty"`
	Tags        []string             `yaml:"tags,omitempty"`
	MidiKey     string               `yaml:"midiKey,omitempty"`
	ControlsMap map[string]m.Control `yaml:"controls"`
	LayersMap   map[string]NewLayer  `yaml:"layers,omitempty"`
}

type NewLayer struct {
	Name        string               `yaml:"name,omitempty" json:"name,omitempty"`
	MidiKey     string               `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	ControlsMap map[string]m.Control `yaml:"controls,omitempty" json:"controls,omitempty"`
}

func TransformKitFormat(path string) error {
	slog.Info("parse dir: " + path)
	items, err := f.ParseYAMLDir(path)
	if err != nil {
		return fmt.Errorf("failed load kit files: %w", err)
	}

	newdir := p.Join(path, "new")
	if err := os.Mkdir(newdir, os.ModePerm); err != nil {
		slog.Error(fmt.Sprintln(fmt.Errorf("failed create dir: %w", err)))
	}

	for k, v := range items {
		switch v.(type) {
		case *m.Instrument:
			break
		case *m.Kit:
			continue
		default:
			continue
		}

		oldi := items[k].(*m.Instrument)
		slog.Info(fmt.Sprintf("process: %s", oldi.Key))

		newi := NewInstrument{
			Id:          oldi.Id,
			Uid:         oldi.Uid,
			Key:         oldi.Key,
			Name:        oldi.Name,
			FullName:    oldi.FullName,
			Type:        oldi.Type,
			SubType:     oldi.SubType,
			Description: oldi.Description,
			Copyright:   oldi.Copyright,
			Licence:     oldi.Licence,
			Credits:     oldi.Credits,
			Tags:        oldi.Tags,
			MidiKey:     oldi.MidiKey,
		}

		// transform instrument controls
		mctrls := trnsformControls(&oldi.Controls)
		newi.ControlsMap = mctrls

		// transform layers
		mlrs := make(map[string]NewLayer, len(oldi.Layers))
		for _, lv := range oldi.Layers {
			lkey := strings.TrimSpace(strings.ToLower(lv.Name))
			var lname string
			if strings.Contains(lkey, " ") {
				lkey = strings.Replace(lkey, " ", "_", -1)
				lname = lv.Name
			}
			nlr := NewLayer{
				MidiKey: lv.MidiKey,
			}
			if len(lname) > 0 {
				nlr.Name = lname
			}
			nlr.ControlsMap = trnsformControls(&lv.Controls)
			mlrs[lkey] = nlr
		}
		newi.LayersMap = mlrs

		// write to file
		data, err := yaml.Marshal(&newi)
		if err != nil {
			slog.Error(fmt.Sprintln(fmt.Errorf("failed marshal to yaml: %w", err)))
		}

		filename := p.Join(newdir, newi.Key+".yaml")
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			slog.Error(fmt.Sprintln(fmt.Errorf("failed writing to file: %s %w", filename, err)))
		}

	}

	return nil
}

func trnsformControls(ctrls *[]m.Control) map[string]m.Control {
	res := make(map[string]m.Control, len(*ctrls))
	for _, cv := range *ctrls {
		ckey := strings.TrimSpace(strings.ToLower(cv.Name))
		var cname string
		if strings.Contains(ckey, " ") {
			ckey = strings.Replace(ckey, " ", "_", -1)
			cname = cv.Name
		}
		ctrl := m.Control{
			Type: cv.Type,
			Key:  cv.Key,
		}
		if len(cname) > 0 {
			ctrl.Name = cname
		}
		res[ckey] = ctrl
	}
	return res
}
