// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// goslicer is the command line program for slicing lines.
//
// Usage
//
// Use `goslicer help` to see full help documentation.
//
//	goslicer [--lines] [--indicies] START[:END] [-|FILE]
//
// Examples
//
//	# show the
//	find . -name '*.go' | goslicer --lines --indicies 0:10
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/goslicer/pkg/slicer"
)

const (
	flagLines    = "lines"
	flagIndicies = "indicies"
)

func initFlags(flag *pflag.FlagSet) {
	flag.BoolP(flagLines, "l", false, "process as lines")
	flag.StringP(flagIndicies, "i", "", "indicies")
}

func stringSlicetoIntSlice(slc []string) ([]int, error) {
	indicies := make([]int, 0)
	for _, str := range slc {
		i, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return indicies, fmt.Errorf("error parsing string %q", str)
		}
		indicies = append(indicies, i)
	}
	return indicies, nil
}

func initReader(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	pathExpanded, err := homedir.Expand(path)
	if err != nil {
		return nil, fmt.Errorf("error expanding input path %q: %w", path, err)
	}
	file, err := os.Open(pathExpanded)
	if err != nil {
		return nil, fmt.Errorf("error opening input file at path %q: %w", path, err)
	}
	return file, nil
}

func main() {

	rootCommand := &cobra.Command{
		Use:                   "goslicer [--lines] [--indicies] START[:END] [-|FILE]",
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    false,
		Short: `goslicer is a simple tool for slicing streams of bytes.
START must be greater than or equal to zero.
END supports negative indicies (as subtracted from the total length).`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()

			if errorBind := v.BindPFlags(cmd.Flags()); errorBind != nil {
				return errorBind
			}

			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			if len(args) > 1 {
				return errors.New("only one input file allowed")
			}

			indiciesString := v.GetString(flagIndicies)
			if len(indiciesString) == 0 {
				return fmt.Errorf("indicies are missing")
			}

			indicies, err := stringSlicetoIntSlice(strings.Split(indiciesString, ":"))
			if err != nil {
				return fmt.Errorf("error parsing indicies: %w", err)
			}

			s, err := slicer.New(false, indicies...)
			if err != nil {
				return fmt.Errorf("error creating slicer: %w", err)
			}

			lines := v.GetBool(flagLines)

			path := "-"
			if len(args) == 1 {
				path = args[0]
			}

			reader, err := initReader(path)
			if err != nil {
				return fmt.Errorf("error initializing reader: %w", err)
			}

			if lines {
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					in := scanner.Bytes()
					out, errSliceBytes := s.SliceBytes(in)
					if errSliceBytes != nil {
						return fmt.Errorf("error slicing bytes %q: %w", in, errSliceBytes)
					}
					if _, errWrite := os.Stdout.Write(append(out, '\n')); errWrite != nil {
						return fmt.Errorf("error writing output to stdout %q: %w", string(out), errWrite)
					}
				}
				if errScanner := scanner.Err(); errScanner != nil {
					return fmt.Errorf("error scanning input: %w", errScanner)
				}
				return nil
			}

			in, err := ioutil.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			out, err := s.SliceBytes(in)
			if err != nil {
				return fmt.Errorf("error slicing bytes %q: %w", in, err)
			}
			if _, err := os.Stdout.Write(append(out, '\n')); err != nil {
				return fmt.Errorf("error writing output to stdout %q", string(out))
			}
			return nil
		},
	}
	initFlags(rootCommand.Flags())

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "goslicer: "+err.Error())
		fmt.Fprintln(os.Stderr, "Try goslicer --help for more information.")
		os.Exit(1)
	}
}
