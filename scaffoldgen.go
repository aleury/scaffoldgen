package scaffoldgen

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	Name            string
	Directory       string
	Repository      string
	HasStaticAssets bool
}

func SetupParseFlags(w io.Writer, args []string) (Config, error) {
	var conf Config
	fs := flag.NewFlagSet("scaffoldgen", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&conf.Name, "n", "", "Project name")
	fs.StringVar(&conf.Directory, "d", "", "Project location on disk")
	fs.StringVar(&conf.Repository, "r", "", "Project remote repository URL")
	fs.BoolVar(&conf.HasStaticAssets, "s", false, "Project will have static assets or not")
	if len(args) == 0 {
		fs.Usage()
		return conf, errors.New("a name, directory, and repository url must be provided")
	}
	err := fs.Parse(args)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func ValidateConfig(conf Config) []error {
	errs := []error{}
	if strings.TrimSpace(conf.Name) == "" {
		errs = append(errs, errors.New("project name cannot be empty"))
	}
	if strings.TrimSpace(conf.Directory) == "" {
		errs = append(errs, errors.New("project directory cannot be empty"))
	}
	if strings.TrimSpace(conf.Repository) == "" {
		errs = append(errs, errors.New("project repository url cannot be empty"))
	}
	return errs
}

func GenerateScaffold(w io.Writer, conf Config) error {
	fmt.Fprintf(w, "Generating %s scaffold at %s...\n", conf.Name, conf.Directory)
	return nil
}

func RunCLI() int {
	conf, err := SetupParseFlags(os.Stdout, os.Args[1:])
	if err != nil {
		return 1
	}
	if errs := ValidateConfig(conf); len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		return 1
	}
	if err := GenerateScaffold(os.Stdout, conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
