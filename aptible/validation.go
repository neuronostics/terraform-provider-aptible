package aptible

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Modified validation.IsURLWithScheme to simply check for a URL that has any scheme and a host.
//
// Ignore linter rule complaining about reimplementing validation.StringNotInSlice.
// lintignore:V013
func validateURL(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" {
		errors = append(errors, fmt.Errorf("expected %q url to not be empty", k))
		return
	}

	u, err := url.Parse(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid url, got %v: %+v", k, v, err))
		return
	}

	if u.Scheme == "" {
		errors = append(errors, fmt.Errorf("expected %q to have a scheme, got %v", k, v))
		return
	}

	if u.Host == "" {
		errors = append(errors, fmt.Errorf("expected %q to have a host, got %v", k, v))
		return
	}

	return
}

// nolint:staticcheck
func errorsToWarnings(validator schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		warns, errs := validator(i, k)
		for _, err := range errs {
			if err != nil {
				warns = append(warns, err.Error())
			}
		}

		return warns, nil
	}
}

var validContainerSizes = []int{
	512,
	1024,
	2048,
	4096,
	7168,
	15360,
	30720,
	61440,
	153600,
	245760,
}
var validateContainerSize = validation.IntInSlice(validContainerSizes)

var validContainerProfiles = []string{
	"m4",
	"r5",
	"c5",
}
var validateContainerProfile = errorsToWarnings(validation.StringInSlice(validContainerProfiles, false))

var validateDiskSize = validation.IntBetween(1, 16000)
