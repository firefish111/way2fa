package creation

import (
	"github.com/firefish111/way2fa/account"
	"fmt"
	"strings"
	"strconv"
)

func validateKey(inp string) error {
	if _, err := account.DecodeTextKey(keyifyKey(inp)); err != nil {
		return fmt.Errorf("2FA Key invalid: %w (got: %s)", err, inp)
	}

	return nil
}

func keyifyKey(inp string) string {
	return strings.ReplaceAll(strings.ToUpper(inp), " ", "")
}

func validateInterv(inp string) error {
	if len(inp) == 0 {
		return nil // to stop atoi from throwing all the toys out the pram
	}

	if _, err := strconv.Atoi(inp); err != nil {
		return fmt.Errorf("Interval, must be an integer number of seconds; %w (got: %s)", err, inp)
	}

	return nil
}

func handlifyAcctId(inp string) string { // just in case user is stupid
	return strings.TrimLeft(inp, "@")
}
