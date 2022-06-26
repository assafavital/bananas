package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

type withCMDFormatter struct {
	delegateFormatter logrus.Formatter
}

func (formatter withCMDFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data[commandLineField] = os.Args
	data, err := formatter.delegateFormatter.Format(entry)
	delete(entry.Data, commandLineField)

	return data, err
}

type stackTraceFormatter struct {
	delegateFormatter logrus.Formatter
}

func errorFormatWithStack(errs []error) string {
	entries := make([]string, len(errs)+1)
	entries[0] = fmt.Sprintf("encountered %d errors:\n-------------------------\n", len(errs))
	for i, err := range errs {
		printedErr := strings.ReplaceAll(fmt.Sprintf("%+v", err), "\n", "\n\t")
		entries[i+1] = fmt.Sprintf("\t* %+v\n", printedErr)
	}

	return strings.Join(entries, "")
}

func formatMultiError(err interface{}) string {
	var errors []error
	switch mErr := err.(type) {
	case *multierror.Error:
		errors = mErr.Errors
	default:
		return ""
	}

	result := ""
	result += errorFormatWithStack(errors)
	for i := range errors {
		result += formatError(errors[i])
	}

	return result
}

func formatError(err interface{}) string {
	// Note: order is somewhat important here, we want to format multi error recursively and print stacktrace for each error
	if multiErrFormat := formatMultiError(err); multiErrFormat != "" {
		return multiErrFormat
	}

	return ""
}

func (f *stackTraceFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	err, ok := entry.Data[logrus.ErrorKey]
	if ok {
		if formatted := formatError(err); formatted != "" {
			entry.Data[stackTraceKey] = formatted
		}
	}

	return f.delegateFormatter.Format(entry)
}
