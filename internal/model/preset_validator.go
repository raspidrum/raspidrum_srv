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
// - other controls of instrument (except `volume` and `pan`) MUST have midiCC
func (p *KitPreset) Validate() error {
	var errs MultiValidationError

	if err := p.indexInstruments(); err != nil {
		errs = append(errs, ValidationError{"preset", err.Error()})
	}

	for _, vc := range p.Channels {
		// validate channel controls
		for _, vcc := range vc.Controls {
			var сve MultiValidationError
			err := vcc.Validate()
			if err != nil {
				if errors.As(err, &сve) {
					errs = append(errs, сve...)
				} else {
					return err
				}
			}
		}

		isManyInstruments := len(vc.instruments) > 1

		for _, vi := range vc.instruments {
			// validate instrument controls
			for _, vic := range vi.Controls {
				var сve MultiValidationError
				err := vic.Validate()
				if err != nil {
					if errors.As(err, &сve) {
						errs = append(errs, сve...)
					} else {
						return err
					}
				}
				if vic.Type != CtrlVolume && vic.Type != CtrlPan && vic.MidiCC == 0 {
					errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s'", vic.Name), "midiCC is required and can't be 0"})
				}
			}

			if len(vi.Layers) == 0 {
				// volume control
				if ctrl, ok := vi.Controls[CtrlVolume]; !ok {
					// many instruments in channel, instrument without layers MUST have `volume` and `pan` controls
					if isManyInstruments {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, CtrlVolume), "is required, but missing"})
					}
				} else {
					if ctrl.MidiCC == 0 {
						// one or many instruments. Instrument without layers. Instrument has control. Control MUST have midiCC
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, CtrlVolume), "midiCC is required and can't be 0"})
					}
				}
				// pan control
				if ctrl, ok := vi.Controls[CtrlPan]; !ok {
					if isManyInstruments {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, CtrlPan), "is required, but missing"})
					}
				} else {
					if ctrl.MidiCC == 0 {
						errs = append(errs, ValidationError{fmt.Sprintf("instrument control '%s.%s'", vi.Name, CtrlPan), "midiCC is required and can't be 0"})
					}
				}
			} else {
				// validate instrument layers
				for kl, vl := range vi.Layers {
					var сve MultiValidationError
					err := vl.Validate()
					if err != nil {
						if errors.As(err, &сve) {
							errs = append(errs, ValidationError{fmt.Sprintf("layer '%s'", kl), сve.Error()})
						} else {
							return err
						}
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
// - controls MAY have `pan` type control.
// - any control MUST have midiCC value
func (p *PresetLayer) Validate() error {
	var errs MultiValidationError
	hasVolume := false

	for k, c := range p.Controls {
		// validate layers controls
		var сve MultiValidationError
		err := c.Validate()
		if err != nil {
			if errors.As(err, &сve) {
				errs = append(errs, сve...)
			} else {
				return err
			}
		}
		cType, ok := ControlTypeFromString[c.Type]
		if !ok {
			continue
		}
		hasVolume = (cType == CTVolume) || hasVolume

		// MIDI CC = 0 - "Bank Select" code. It can't be used for layer control
		if c.MidiCC == 0 {
			errs = append(errs, ValidationError{fmt.Sprintf("control '%s'", k), "midiCC is required and can't be 0"})
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
// - value MUST be in range 0..127 if MidiCC is not 0
// - value MUST be in range -1.0..1.0 for pan if MidiCC is 0
// - value MUST be in range 0..1.0 for volume if MidiCC is 0
func (c *PresetControl) Validate() error {
	var errs MultiValidationError
	if _, ok := ControlTypeFromString[c.Type]; !ok {
		errs = append(errs, ValidationError{"type", fmt.Sprintf("unknown value '%s'", c.Type)})
	}
	if c.MidiCC != 0 && (c.Value < 0 || c.Value > 127) {
		errs = append(errs, ValidationError{"value", "value must be in range 0..127"})
	}
	if c.MidiCC == 0 && (c.Type == CtrlPan && (c.Value < -1.0 || c.Value > 1.0)) {
		errs = append(errs, ValidationError{"value", "value must be in range -1.0..1.0"})
	}
	if c.MidiCC == 0 && (c.Type == CtrlVolume && (c.Value < 0 || c.Value > 1.0)) {
		errs = append(errs, ValidationError{"value", "value must be in range 0..1.0"})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
