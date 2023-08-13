package scaffoldgen_test

import (
	"bytes"
	"errors"
	"scaffoldgen"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetupParseFlags_DisplayUsageWhenGivenNoArgs(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	args := []string{}
	want := "Usage of scaffoldgen:"
	wantErr := "a name, directory, and repository url must be provided"
	_, err := scaffoldgen.SetupParseFlags(buf, args)
	if !cmp.Equal(wantErr, err.Error()) {
		t.Error(cmp.Diff(wantErr, err.Error()))
	}
	got := buf.String()
	if !strings.Contains(got, want) {
		t.Errorf("expected usage to contain: %s", want)
	}
}

func TestSetupParseFlags_ReturnsValidConfig(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	args := []string{"-d", "./project1", "-n", "project1", "-r", "github.com/username/project1"}
	want := scaffoldgen.Config{
		Name:            "project1",
		Directory:       "./project1",
		Repository:      "github.com/username/project1",
		HasStaticAssets: false,
	}
	got, err := scaffoldgen.SetupParseFlags(buf, args)
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSetupParseFlags_ReturnsErrorForHelpFlag(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	want := errors.New("flag: help requested")
	_, err := scaffoldgen.SetupParseFlags(buf, []string{"-h"})
	if !cmp.Equal(want.Error(), err.Error()) {
		t.Error(cmp.Diff(want.Error(), err.Error()))
	}
}

func TestValidateConfig_ReturnsZeroErrorsWhenGivenValidInput(t *testing.T) {
	t.Parallel()
	conf := scaffoldgen.Config{
		Name:            "project1",
		Directory:       "./project1",
		Repository:      "github.com/username/project1",
		HasStaticAssets: true,
	}
	want := []error{}
	got := scaffoldgen.ValidateConfig(conf)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestValidateConfig_ReturnsErrorsWhenMissingRequiredParameters(t *testing.T) {
	t.Parallel()
	conf := scaffoldgen.Config{}
	want := []string{
		"project name cannot be empty",
		"project directory cannot be empty",
		"project repository url cannot be empty",
	}
	got := mapErrorsToStrings(t, scaffoldgen.ValidateConfig(conf))
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGenerateScaffold(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	conf := scaffoldgen.Config{
		Name:            "project1",
		Directory:       "./project1",
		Repository:      "github.com/username/project1",
		HasStaticAssets: false,
	}
	err := scaffoldgen.GenerateScaffold(buf, conf)
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	want := "Generating project1 scaffold at ./project1...\n"
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func mapErrorsToStrings(t testing.TB, errs []error) []string {
	t.Helper()
	result := []string{}
	for _, err := range errs {
		result = append(result, err.Error())
	}
	return result
}
