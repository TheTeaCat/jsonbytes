package jsonbytes

import (
	"errors"
	"fmt"
)

type jsonValidator struct {
	json       []byte
	jsonLength int
	readIndex  int
	readHead   byte
}

func newJsonValidator(json []byte) (*jsonValidator, error) {
	if len(json) == 0 {
		return nil, errors.New("jsonvalidator needs more than zero bytes")
	}
	return &jsonValidator{
		json:       json,
		readHead:   json[0],
		jsonLength: len(json),
		readIndex:  0,
	}, nil
}

func (state *jsonValidator) consumeValue() error {
	state.consumeWhitespace()
	if state.readIndex == state.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	switch state.readHead {
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
		return state.errorUnexpectedCharacter("any of \"10123456789{[tfn")
	}
	if err != nil {
		return err
	}
	state.consumeWhitespace()
	return nil
}

func (state *jsonValidator) consumeObject() error {
	state.readUnsafe()
	state.consumeWhitespace()
	if state.readIndex == state.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	for {
		switch state.readHead {
		case ',':
			err = state.consumeByte(',')
			if err != nil {
				return err
			}
		case '}':
			return state.consumeByte('}')
		default:
			state.consumeWhitespace()
			if state.readIndex == state.jsonLength {
				return errors.New("read head ran out of json")
			}
			err = state.consumeName()
			if err != nil {
				return err
			}
			state.consumeWhitespace()
			if state.readIndex == state.jsonLength {
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

func (state *jsonValidator) consumeArray() error {
	state.readUnsafe()
	state.consumeWhitespace()
	if state.readIndex == state.jsonLength {
		return errors.New("read head ran out of json")
	}
	var err error
	for {
		switch state.readHead {
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

func (state *jsonValidator) consumeWhitespace() {
	for state.readIndex < state.jsonLength && (state.readHead == ' ' ||
		state.readHead == '\t' ||
		state.readHead == '\n' ||
		state.readHead == '\r') {
		state.readUnsafe()
	}
}

func (state *jsonValidator) consumeString() error {
	state.readUnsafe()
	prevHead := byte(0)
	for (state.readHead != '"' || prevHead == '\\') && state.readIndex < state.jsonLength {
		if state.readHead < 32 || state.readHead == 127 {
			return state.errorUnexpectedCharacter("any codepoint except \" or \\ or control characters")
		}
		if prevHead == '\\' {
			switch state.readHead {
			case '"', '/', '\\', 'b', 'f', 'n', 'r', 't':
				break
			case 'u':
				for i := 0; i < 4 && state.readIndex < state.jsonLength; i++ {
					state.readUnsafe()
					if state.readIndex == state.jsonLength {
						return state.errorUnexpectedCharacter("4 hex digits")
					}
					if !(('0' <= state.readHead && state.readHead <= '9') ||
						('A' <= state.readHead && state.readHead <= 'F') ||
						('a' <= state.readHead && state.readHead <= 'f')) {
						return state.errorUnexpectedCharacter("4 hex digits")
					}
				}
			default:
				return state.errorUnexpectedCharacter("any of \"/\\bfnrtu")
			}

		}
		prevHead = state.readHead
		state.readUnsafe()
	}
	return state.consumeByte('"')
}

func (state *jsonValidator) consumeName() error {
	state.readUnsafe()
	prevHead := byte(0)
	for state.readHead != '"' || prevHead == '\\' {
		prevHead = state.readHead
		state.readIndex += 1
		if state.readIndex > state.jsonLength {
			return state.errorUnexpectedCharacter("\"")
		} else if state.readIndex != state.jsonLength {
			state.readHead = state.json[state.readIndex]
		}
	}
	return state.consumeByte('"')
}

func (state *jsonValidator) consumeNumber() error {
	if state.readHead == '-' {
		state.readUnsafe()
	}
	if state.readHead == '0' {
		state.readUnsafe()
	} else if 49 <= state.readHead && state.readHead <= 57 {
		for 48 <= state.readHead && state.readHead <= 57 && state.readIndex < state.jsonLength {
			state.readUnsafe()
		}
	} else {
		return state.errorUnexpectedCharacter("any of 0123456789")
	}
	if state.readHead == '.' {
		state.readUnsafe()
		if 48 > state.readHead || state.readHead > 57 {
			return state.errorUnexpectedCharacter("any of 0123456789")
		}
		for 48 <= state.readHead && state.readHead <= 57 && state.readIndex < state.jsonLength {
			state.readUnsafe()
		}
	}
	if state.readHead == 'E' || state.readHead == 'e' {
		state.readUnsafe()
		if state.readHead != '+' && state.readHead != '-' {
			return state.errorUnexpectedCharacter("+ or -")
		}
		state.readUnsafe()
		if state.readIndex == state.jsonLength {
			return errors.New("read head ran out of json")
		}
		for 48 <= state.readHead && state.readHead <= 57 && state.readIndex < state.jsonLength {
			state.readUnsafe()
		}
	}
	return nil
}

func (state *jsonValidator) consumeTrue() error {
	state.readUnsafe()
	return state.consumeSlice([]byte("rue"))
}

func (state *jsonValidator) consumeFalse() error {
	state.readUnsafe()
	for i := 0; i < 4; i++ {
		if state.readHead != []byte("alse")[i] {
			return state.errorUnexpectedCharacter(string([]byte("alse")[i]))
		}
		state.readUnsafe()
	}
	return nil
}

func (state *jsonValidator) consumeNull() error {
	state.readUnsafe()
	return state.consumeSlice([]byte("ull"))
}

func (state *jsonValidator) consumeSlice(slice []byte) error {
	for _, expectedByte := range slice {
		err := state.consumeByte(expectedByte)
		if err != nil {
			return err
		}
	}
	return nil
}

func (state *jsonValidator) consumeByte(expectedByte byte) error {
	if state.readHead != expectedByte {
		return state.errorUnexpectedCharacter(string(expectedByte))
	}
	state.readIndex += 1
	if state.readIndex > state.jsonLength {
		return errors.New("read head ran out of json")
	} else if state.readIndex != state.jsonLength {
		state.readHead = state.json[state.readIndex]
	}
	return nil
}

func (state *jsonValidator) readUnsafe() {
	state.readIndex += 1
	if state.readIndex != state.jsonLength {
		state.readHead = state.json[state.readIndex]
	}
}

func (state *jsonValidator) errorUnexpectedCharacter(expectedBytes string) error {
	if state.readIndex >= state.jsonLength {
		return fmt.Errorf("expected %s but reached end of json", expectedBytes)
	}
	return fmt.Errorf(
		"expected %s at index %d but read '%s'",
		expectedBytes, state.readIndex, string(state.readHead),
	)
}
