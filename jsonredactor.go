package jsonbytes

import "errors"

type jsonRedactor struct {
	jsonValidator *jsonValidator
	writeIndex    int
}

func newJsonRedactor(json []byte) (*jsonRedactor, error) {
	jsonValidator, err := newJsonValidator(json)
	if err != nil {
		return nil, err
	}
	return &jsonRedactor{
		jsonValidator: jsonValidator,
		writeIndex:    0,
	}, nil
}

func (state *jsonRedactor) consumeValue() error {
	state.jsonValidator.consumeWhitespace()
	if state.jsonValidator.readIndex == state.jsonValidator.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	switch state.jsonValidator.readHead {
	case '"':
		err = state.consumeString()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err = state.consumeNumber()
	case '{':
		err = state.consumeObject()
	case '[':
		err = state.consumeArray()
	case 't':
		err = state.consumeTrue()
	case 'f':
		err = state.consumeFalse()
	case 'n':
		err = state.consumeNull()
	default:
		return state.jsonValidator.errorUnexpectedCharacter("any of \"10123456789{[tfn")
	}
	if err != nil {
		return err
	}
	state.jsonValidator.consumeWhitespace()
	return nil
}

func (state *jsonRedactor) consumeObject() error {
	state.writeUnsafe()
	state.jsonValidator.consumeWhitespace()
	if state.jsonValidator.readIndex == state.jsonValidator.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	for {
		switch state.jsonValidator.readHead {
		case ',':
			err = state.consumeByte(',')
			if err != nil {
				return err
			}
		case '}':
			return state.consumeByte('}')
		default:
			state.jsonValidator.consumeWhitespace()
			if state.jsonValidator.readIndex == state.jsonValidator.jsonLength {
				return errors.New("read head ran out of json")
			}
			err = state.consumeName()
			if err != nil {
				return err
			}
			state.jsonValidator.consumeWhitespace()
			if state.jsonValidator.readIndex == state.jsonValidator.jsonLength {
				return errors.New("read head ran out of json")
			}
			err = state.consumeByte(':')
			if err != nil {
				return err
			}
			err = state.consumeValue()
			if err != nil {
				return err
			}
		}
	}
}

func (state *jsonRedactor) consumeArray() error {
	state.writeUnsafe()
	state.jsonValidator.consumeWhitespace()
	if state.jsonValidator.readIndex == state.jsonValidator.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	for {
		switch state.jsonValidator.readHead {
		case ']':
			return state.consumeByte(']')
		case ',':
			err = state.consumeByte(',')
			if err != nil {
				return err
			}
		default:
			err = state.consumeValue()
			if err != nil {
				return err
			}
		}
	}
}

func (state *jsonRedactor) consumeString() error {
	err := state.jsonValidator.consumeString()
	if err != nil {
		return err
	}
	state.jsonValidator.json[state.writeIndex] = '"'
	state.writeIndex += 1
	state.jsonValidator.json[state.writeIndex] = '"'
	state.writeIndex += 1
	return nil
}

func (state *jsonRedactor) consumeName() error {
	bytesConsumed := -state.jsonValidator.readIndex
	err := state.jsonValidator.consumeString()
	if err != nil {
		return err
	}
	bytesConsumed += state.jsonValidator.readIndex
	copy(
		state.jsonValidator.json[state.writeIndex:state.writeIndex+bytesConsumed],
		state.jsonValidator.json[state.jsonValidator.readIndex-bytesConsumed:state.jsonValidator.readIndex],
	)
	state.writeIndex += bytesConsumed
	return nil
}

func (state *jsonRedactor) consumeNumber() error {
	err := state.jsonValidator.consumeNumber()
	if err != nil {
		return err
	}
	state.jsonValidator.json[state.writeIndex] = '0'
	state.writeIndex += 1
	return nil
}

func (state *jsonRedactor) consumeTrue() error {
	err := state.jsonValidator.consumeTrue()
	if err != nil {
		return err
	}
	copy(state.jsonValidator.json[state.writeIndex:state.writeIndex+4], []byte("true"))
	state.writeIndex += 4
	return nil
}

func (state *jsonRedactor) consumeFalse() error {
	err := state.jsonValidator.consumeFalse()
	if err != nil {
		return err
	}
	copy(state.jsonValidator.json[state.writeIndex:state.writeIndex+4], []byte("true"))
	state.writeIndex += 4
	return nil
}

func (state *jsonRedactor) consumeNull() error {
	err := state.jsonValidator.consumeNull()
	if err != nil {
		return err
	}
	copy(state.jsonValidator.json[state.writeIndex:state.writeIndex+4], []byte("null"))
	state.writeIndex += 4
	return nil
}

func (state *jsonRedactor) consumeByte(expectedByte byte) error {
	err := state.jsonValidator.consumeByte(expectedByte)
	if err != nil {
		return err
	}
	state.jsonValidator.json[state.writeIndex] = expectedByte
	state.writeIndex += 1
	return nil
}

func (state *jsonRedactor) writeUnsafe() {
	if state.writeIndex != state.jsonValidator.readIndex {
		state.jsonValidator.json[state.writeIndex] = state.jsonValidator.readHead
	}
	state.writeIndex += 1
	state.jsonValidator.readUnsafe()
}
