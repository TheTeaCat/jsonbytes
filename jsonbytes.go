package jsonbytes

import "errors"

// IsJson takes a single argument maybeJson []byte and returns nil if maybeJson is a valid JSON value, else an error
// detailing why it is not a valid JSON value.
func IsJson(maybeJson []byte) error {
	jsonValidator, err := newJsonValidator(maybeJson)
	if err != nil {
		return err
	}
	err = jsonValidator.consumeValue()
	if err != nil {
		return err
	}
	if jsonValidator.readIndex != jsonValidator.jsonLength {
		return errors.New("failed to consume entire json string")
	}
	return nil
}

// RedactAllValues takes a single argument inputJson []byte and returns a new []byte which will be identical to
// inputJson with all string values replaced with "", numbers replaced with 0, booleans replaced with true and
// unecessary whitespace characters removed. If inputJson is not a valid JSON value, RedactAllValues will return an
// error explaining why.
func RedactAllValues(inputJson []byte) ([]byte, error) {
	jsonRedactor, err := newJsonRedactor(inputJson)
	if err != nil {
		return nil, err
	}
	err = jsonRedactor.consumeValue()
	if err != nil {
		return nil, err
	}
	if jsonRedactor.jsonValidator.readIndex != jsonRedactor.jsonValidator.jsonLength {
		return nil, errors.New("failed to consume entire json string")
	}
	return jsonRedactor.jsonValidator.json[:jsonRedactor.writeIndex], nil
}
