package creation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/firefish111/way2fa/account"
)

// ensure key is valid base64, and trim spaces (as most often services give them in space-separated chunks)
func keyifyKey(inp string) string {
	return strings.ReplaceAll(strings.ToUpper(inp), " ", "")
}

// trim initial @ signs from account id, as those arer implicit and added later in UI only
func handlifyAcctId(inp string) string {
	return strings.TrimLeft(inp, "@")
}

func validateKey(inp string) error {
	if _, err := account.DecodeTextKey(keyifyKey(inp)); err != nil {
		return fmt.Errorf("2FA Key invalid: %w (got: %s)", err, inp)
	}

	return nil
}

func validateInterv(inp string) error {
	if len(inp) == 0 { // if no string given, is valid, as is optional
		return nil // to stop atoi from throwing all the toys out the pram
	}

	n, err := strconv.Atoi(inp)
	if err != nil {
		return fmt.Errorf("Interval must be an integer number of seconds; %w (got: %s)", err, inp)
	}

	if n == 0 {
		return fmt.Errorf("Interval muust be a positive integer (got %d)", n)
	}

	return nil
}

func validateService(inp string) error {
	if len(inp) == 0 {
		return fmt.Errorf("Must supply service name")
	}

	return nil
}
