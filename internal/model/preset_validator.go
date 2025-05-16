package model

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type MultiValidationError []ValidationError

func (mve MultiValidationError) Error() string {
	var msgs []string
	for _, err := range mve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Validations:
// - in case many instruments in channel, instrument without layers MUST have `volume` and `pan` controls. Its controls MUST have midiCC
// - in case one instrument in channel,   instrument without layers MAY  have `volume` or `pan` controls.  Its controls MUST have midiCC
// - in case one or many instrument in channel, instrument with layers MAY have `volume` or `pan` controls. Its controls MAY have midiCC
func (p *KitPreset) Validate() error {
	var errs MultiValidationError
	ctrlVolumeKey := ControlTypeToString[CTVolume]
	ctrlPanKey := ControlTypeToString[CTPan]

	if err := p.indexInstruments(); err != nil {
		errs = append(errs, ValidationError{"preset", err.Error()})
	}

	for _, vc := range p.Channels {
		// validate channel controls
		for _, vcc := range vc.Controls {
			var сve MultiValidationError
			err := vcc.Validate()
			if errors.As(err, &сve) {
				errs = append(errs, сve...)
			} else {
				return err
			}
		}

		manyInstruments := len(vc.instruments) > 1

		for _, vi := range vc.instruments {
			// validate instrument controls
			for _, vic := range vi.Controls {
				var сve MultiValidationError
				err := vic.Validate()
				if errors.As(err, &сve) {
					errs = append(errs, сve...)
				} else {
					return err
				}
			}

			if len(vi.Layers) == 0 {
				// volume control
				if ctrl, ok := vi.Controls[ctrlVolumeKey]; !ok {
					// many instruments in channel, instrument without layers MUST have `volume` and `pan` controls
					if manyInstruments {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, ctrlVolumeKey), "is required, but missing"})
					}
				} else {
					if ctrl.MidiCC == 0 {
						// one or many instruments. Instrument without layers. Instrument has control. Control MUST have midiCC
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, ctrlVolumeKey), "is required and can't be 0"})
					}
				}
				// pan control
				if ctrl, ok := vi.Controls[ctrlPanKey]; !ok {
					if manyInstruments {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, ctrlPanKey), "is required, but missing"})
					}
				} else {
					if ctrl.MidiCC == 0 {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, ctrlPanKey), "is required and can't be 0"})
					}
				}
			} else {
				// validate instrument layers
				for _, vl := range vi.Layers {
					var сve MultiValidationError
					err := vl.Validate()
					if errors.As(err, &сve) {
						errs = append(errs, сve...)
					} else {
						return err
					}
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Validations:
// - controls MUST have `volume` type control. It control MUST have midiCC
// - controls MAY have `pan` type control. If `pan` exists then it MUST have midiCC value
func (p *PresetLayer) Validate() error {
	var errs MultiValidationError
	hasVolume := false

	for k, c := range p.Controls {
		// validate layers controls
		var сve MultiValidationError
		err := c.Validate()
		if errors.As(err, &сve) {
			errs = append(errs, сve...)
		} else {
			return err
		}

		cType, ok := ControlTypeFromString[c.Type]
		if !ok {
			continue
		}
		hasVolume = (cType == CTVolume) || hasVolume

		if cType == CTVolume || cType == CTPan {
			// MIDI CC = 0 - "Bank Select" code. It can't be used for layer control
			if c.MidiCC == 0 {
				errs = append(errs, ValidationError{fmt.Sprintf("control '%s'", k), "is required and can't be 0"})
			}
		}
	}

	if !hasVolume {
		errs = append(errs, ValidationError{" ", "missing 'volume' control"})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Validations:
// - type MUST match one of ControlType
func (c *PresetControl) Validate() error {
	var errs MultiValidationError
	if _, ok := ControlTypeFromString[c.Type]; !ok {
		errs = append(errs, ValidationError{"type", fmt.Sprintf("unknown value '%s'", c.Type)})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
