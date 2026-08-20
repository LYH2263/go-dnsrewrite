package rule

import "strings"

// Validate 校验内部规则。
func Validate(r Rule) error {
	if strings.TrimSpace(r.ID) == "" {
		return errInvalid
	}
	if strings.TrimSpace(r.Pattern) == "" {
		return errInvalid
	}
	switch r.Kind {
	case KindExact, KindSuffix, KindGlob, "":
	default:
		return errInvalid
	}
	switch r.Action {
	case ActionRewrite, ActionForward, ActionRefuse, "":
	default:
		return errInvalid
	}
	return nil
}

type ruleError string

func (e ruleError) Error() string { return string(e) }

const errInvalid ruleError = "rule: invalid"
