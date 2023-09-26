// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/extensions/filters/network/thrift_proxy/filters/payload_to_metadata/v3/payload_to_metadata.proto

package payload_to_metadatav3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on PayloadToMetadata with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PayloadToMetadata) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayloadToMetadata with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PayloadToMetadataMultiError, or nil if none found.
func (m *PayloadToMetadata) ValidateAll() error {
	return m.validate(true)
}

func (m *PayloadToMetadata) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetRequestRules()) < 1 {
		err := PayloadToMetadataValidationError{
			field:  "RequestRules",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetRequestRules() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PayloadToMetadataValidationError{
						field:  fmt.Sprintf("RequestRules[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PayloadToMetadataValidationError{
						field:  fmt.Sprintf("RequestRules[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PayloadToMetadataValidationError{
					field:  fmt.Sprintf("RequestRules[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PayloadToMetadataMultiError(errors)
	}
	return nil
}

// PayloadToMetadataMultiError is an error wrapping multiple validation errors
// returned by PayloadToMetadata.ValidateAll() if the designated constraints
// aren't met.
type PayloadToMetadataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayloadToMetadataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayloadToMetadataMultiError) AllErrors() []error { return m }

// PayloadToMetadataValidationError is the validation error returned by
// PayloadToMetadata.Validate if the designated constraints aren't met.
type PayloadToMetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayloadToMetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayloadToMetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayloadToMetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayloadToMetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayloadToMetadataValidationError) ErrorName() string {
	return "PayloadToMetadataValidationError"
}

// Error satisfies the builtin error interface
func (e PayloadToMetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayloadToMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayloadToMetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayloadToMetadataValidationError{}

// Validate checks the field values on PayloadToMetadata_KeyValuePair with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PayloadToMetadata_KeyValuePair) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayloadToMetadata_KeyValuePair with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// PayloadToMetadata_KeyValuePairMultiError, or nil if none found.
func (m *PayloadToMetadata_KeyValuePair) ValidateAll() error {
	return m.validate(true)
}

func (m *PayloadToMetadata_KeyValuePair) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for MetadataNamespace

	if utf8.RuneCountInString(m.GetKey()) < 1 {
		err := PayloadToMetadata_KeyValuePairValidationError{
			field:  "Key",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := PayloadToMetadata_ValueType_name[int32(m.GetType())]; !ok {
		err := PayloadToMetadata_KeyValuePairValidationError{
			field:  "Type",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	switch m.ValueType.(type) {

	case *PayloadToMetadata_KeyValuePair_Value:
		// no validation rules for Value

	case *PayloadToMetadata_KeyValuePair_RegexValueRewrite:

		if all {
			switch v := interface{}(m.GetRegexValueRewrite()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PayloadToMetadata_KeyValuePairValidationError{
						field:  "RegexValueRewrite",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PayloadToMetadata_KeyValuePairValidationError{
						field:  "RegexValueRewrite",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetRegexValueRewrite()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PayloadToMetadata_KeyValuePairValidationError{
					field:  "RegexValueRewrite",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PayloadToMetadata_KeyValuePairMultiError(errors)
	}
	return nil
}

// PayloadToMetadata_KeyValuePairMultiError is an error wrapping multiple
// validation errors returned by PayloadToMetadata_KeyValuePair.ValidateAll()
// if the designated constraints aren't met.
type PayloadToMetadata_KeyValuePairMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayloadToMetadata_KeyValuePairMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayloadToMetadata_KeyValuePairMultiError) AllErrors() []error { return m }

// PayloadToMetadata_KeyValuePairValidationError is the validation error
// returned by PayloadToMetadata_KeyValuePair.Validate if the designated
// constraints aren't met.
type PayloadToMetadata_KeyValuePairValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayloadToMetadata_KeyValuePairValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayloadToMetadata_KeyValuePairValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayloadToMetadata_KeyValuePairValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayloadToMetadata_KeyValuePairValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayloadToMetadata_KeyValuePairValidationError) ErrorName() string {
	return "PayloadToMetadata_KeyValuePairValidationError"
}

// Error satisfies the builtin error interface
func (e PayloadToMetadata_KeyValuePairValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayloadToMetadata_KeyValuePair.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayloadToMetadata_KeyValuePairValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayloadToMetadata_KeyValuePairValidationError{}

// Validate checks the field values on PayloadToMetadata_Rule with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PayloadToMetadata_Rule) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayloadToMetadata_Rule with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PayloadToMetadata_RuleMultiError, or nil if none found.
func (m *PayloadToMetadata_Rule) ValidateAll() error {
	return m.validate(true)
}

func (m *PayloadToMetadata_Rule) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetFieldSelector() == nil {
		err := PayloadToMetadata_RuleValidationError{
			field:  "FieldSelector",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetFieldSelector()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "FieldSelector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "FieldSelector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFieldSelector()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PayloadToMetadata_RuleValidationError{
				field:  "FieldSelector",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetOnPresent()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "OnPresent",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "OnPresent",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOnPresent()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PayloadToMetadata_RuleValidationError{
				field:  "OnPresent",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetOnMissing()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "OnMissing",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PayloadToMetadata_RuleValidationError{
					field:  "OnMissing",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOnMissing()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PayloadToMetadata_RuleValidationError{
				field:  "OnMissing",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.MatchSpecifier.(type) {

	case *PayloadToMetadata_Rule_MethodName:
		// no validation rules for MethodName

	case *PayloadToMetadata_Rule_ServiceName:
		// no validation rules for ServiceName

	default:
		err := PayloadToMetadata_RuleValidationError{
			field:  "MatchSpecifier",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if len(errors) > 0 {
		return PayloadToMetadata_RuleMultiError(errors)
	}
	return nil
}

// PayloadToMetadata_RuleMultiError is an error wrapping multiple validation
// errors returned by PayloadToMetadata_Rule.ValidateAll() if the designated
// constraints aren't met.
type PayloadToMetadata_RuleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayloadToMetadata_RuleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayloadToMetadata_RuleMultiError) AllErrors() []error { return m }

// PayloadToMetadata_RuleValidationError is the validation error returned by
// PayloadToMetadata_Rule.Validate if the designated constraints aren't met.
type PayloadToMetadata_RuleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayloadToMetadata_RuleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayloadToMetadata_RuleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayloadToMetadata_RuleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayloadToMetadata_RuleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayloadToMetadata_RuleValidationError) ErrorName() string {
	return "PayloadToMetadata_RuleValidationError"
}

// Error satisfies the builtin error interface
func (e PayloadToMetadata_RuleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayloadToMetadata_Rule.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayloadToMetadata_RuleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayloadToMetadata_RuleValidationError{}

// Validate checks the field values on PayloadToMetadata_FieldSelector with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PayloadToMetadata_FieldSelector) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PayloadToMetadata_FieldSelector with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// PayloadToMetadata_FieldSelectorMultiError, or nil if none found.
func (m *PayloadToMetadata_FieldSelector) ValidateAll() error {
	return m.validate(true)
}

func (m *PayloadToMetadata_FieldSelector) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := PayloadToMetadata_FieldSelectorValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if val := m.GetId(); val < -32768 || val > 32767 {
		err := PayloadToMetadata_FieldSelectorValidationError{
			field:  "Id",
			reason: "value must be inside range [-32768, 32767]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetChild()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PayloadToMetadata_FieldSelectorValidationError{
					field:  "Child",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PayloadToMetadata_FieldSelectorValidationError{
					field:  "Child",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetChild()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PayloadToMetadata_FieldSelectorValidationError{
				field:  "Child",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PayloadToMetadata_FieldSelectorMultiError(errors)
	}
	return nil
}

// PayloadToMetadata_FieldSelectorMultiError is an error wrapping multiple
// validation errors returned by PayloadToMetadata_FieldSelector.ValidateAll()
// if the designated constraints aren't met.
type PayloadToMetadata_FieldSelectorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PayloadToMetadata_FieldSelectorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PayloadToMetadata_FieldSelectorMultiError) AllErrors() []error { return m }

// PayloadToMetadata_FieldSelectorValidationError is the validation error
// returned by PayloadToMetadata_FieldSelector.Validate if the designated
// constraints aren't met.
type PayloadToMetadata_FieldSelectorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PayloadToMetadata_FieldSelectorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PayloadToMetadata_FieldSelectorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PayloadToMetadata_FieldSelectorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PayloadToMetadata_FieldSelectorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PayloadToMetadata_FieldSelectorValidationError) ErrorName() string {
	return "PayloadToMetadata_FieldSelectorValidationError"
}

// Error satisfies the builtin error interface
func (e PayloadToMetadata_FieldSelectorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPayloadToMetadata_FieldSelector.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PayloadToMetadata_FieldSelectorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PayloadToMetadata_FieldSelectorValidationError{}
