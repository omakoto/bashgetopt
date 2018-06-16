package main

import (
	"errors"
	"go/types"
	"strings"
)

type OptionSpec struct {
	longOption  string
	shortOption rune

	evalString string
	help       string

	// only supports IsBoolean, IsInteger and IsString.
	argType types.BasicInfo

	boolValue bool
	// intValue    int // For integer we use stringValue and parse by ourselves.
	stringValue string
}

func parseFlagSpec(specStr string) (*OptionSpec, error) {
	spec := OptionSpec{}

	specStr = strings.TrimSpace(specStr)

	// Extract the help.
	ar := strings.SplitN(specStr, "#", 2)

	if len(ar) > 1 {
		spec.help = strings.TrimSpace(ar[1])
	}

	// Split flags and command
	ar = strings.SplitN(ar[0], " ", 2)

	if len(ar) < 2 {
		return nil, errors.New("Invalid option spec; missing eval string: " + specStr)
	}
	spec.evalString = strings.TrimSpace(ar[1])

	// Split flags
	flagSpec := ar[0]

	// Detect the type.
	spec.argType = types.IsBoolean
	if strings.HasSuffix(flagSpec, ":") {
		flagSpec = flagSpec[0 : len(flagSpec)-1]
		spec.argType = types.IsString
	} else if strings.HasSuffix(flagSpec, "=s") {
		flagSpec = flagSpec[0 : len(flagSpec)-2]
		spec.argType = types.IsString
	} else if strings.HasSuffix(flagSpec, "=i") {
		flagSpec = flagSpec[0 : len(flagSpec)-2]
		spec.argType = types.IsInteger
	}

	if len(flagSpec) == 0 {
		return nil, errors.New("Invalid option spec; empty flag spec: " + specStr)
	}

	flags := strings.Split(flagSpec, "|")
	if len(flags) > 2 {
		return nil, errors.New("Invalid option spec; too many flags: " + specStr)
	}

	for _, flag := range flags {
		if len(flag) == 1 {
			// Short
			if spec.shortOption != 0 {
				return nil, errors.New("Invalid option spec; multiple short flags: " + specStr)
			}
			spec.shortOption = rune(flag[0])
		} else {
			// Long
			if spec.longOption != "" {
				return nil, errors.New("Invalid option spec; multiple long flags: " + specStr)
			}
			spec.longOption = flag
		}
	}
	return &spec, nil
}

func ParseOptionSpec(specStr string) ([]*OptionSpec, error) {
	specList := strings.Split(specStr, "\n")

	ret := make([]*OptionSpec, 0)
	for _, specStr := range specList {
		specStr = strings.TrimSpace(specStr)
		if len(specStr) == 0 {
			continue
		}

		spec, err := parseFlagSpec(specStr)
		if err != nil {
			return nil, err
		}
		ret = append(ret, spec)
	}

	return ret, nil
}
