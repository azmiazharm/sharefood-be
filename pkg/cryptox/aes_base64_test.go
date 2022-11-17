// Package cryptox
package cryptox

import (
	"crypto/aes"
	"testing"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestEncryptBase64AES(t *testing.T) {

	var testCase = []struct {
		Plain string
		Key   string
		Error error
	}{
		{
			Plain: `{"id":12345678901234567000,"client_id":12345,"client_reference_id":"abc11313134100001","product_code":"PLN","customer_id":"8212919291921"}`,
			Key:   "123",
			Error: aes.KeySizeError(3),
		},
		{
			Plain: `{"id":12345678901234567000,"client_id":12345,"client_reference_id":"abc11313134100001","product_code":"PLN","customer_id":"8212919291921"}`,
			Key:   "kXrtDw$3RhkX8g@#",
			Error: nil,
		},
		{
			Plain: ``,
			Key:   "kXrtDw$3RhkX8g@#",
			Error: nil,
		},
	}

	for _, x := range testCase {

		_, err := EncryptBase64AES(x.Plain, []byte(x.Key))

		if err == x.Error {
			t.Logf("%s expected %v, got %v", success, x.Error, err)
		} else {
			t.Errorf("%s expected %v, got %v", failed, x.Error, err)
		}
	}

}

func TestDecryptBase64AES2(t *testing.T) {
	var testCase = []struct {
		Plain       string
		ChipperText string
		Key         string
		Error       error
	}{
		{
			Plain:       ``,
			Key:         "123",
			ChipperText: ``,
			Error:       aes.KeySizeError(3),
		},
		{
			Plain:       `test-case`,
			Key:         "1234567890123456",
			ChipperText: `HZkiAY0Q2YhO0uZXOsHewoPDlvYcdRCytA==`,
			Error:       nil,
		},
	}

	for _, x := range testCase {

		result, err := DecryptBase64AES(x.ChipperText, []byte(x.Key))

		if err == x.Error {
			t.Logf("%s expected %v, got %v", success, x.Error, err)
		} else {
			t.Errorf("%s expected %v, got %v", failed, x.Error, err)
		}

		if result == x.Plain {
			t.Logf("%s expected %v, got %v", success, x.Plain, result)
		} else {
			t.Errorf("%s expected %v, got %v", failed, x.Plain, result)
		}
	}

}
