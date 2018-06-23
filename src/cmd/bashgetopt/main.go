package main

import (
	"bytes"
	"fmt"
	"github.com/omakoto/compromise/src/compromise/compmain"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/pborman/getopt/v2"
	"go/types"
	"os"
	"strconv"
	"strings"
)

const (
	bashCompletionFlag            = "bash-completion"
	realBashCompletionCommandFlag = "bash-completion-command"
)

var (
	description  = getopt.StringLong("description", 'd', "", "Specify the command description")
	usage        = getopt.StringLong("usage", 'u', "", "Specify function name that prints usage")
	_            = getopt.BoolLong("allow-files", 'F', "Build command completion that allows files (default)")
	noAllowFiles = getopt.BoolLong("no-allow-files", 'N', "Build command completion that doesn't allow files")
	_            = getopt.BoolLong("", 'x', "(Unused, kept for backward compatibility)")

	bashCompletionCommand = getopt.StringLong(realBashCompletionCommandFlag, 0, "")

	dump = os.Getenv("BASHGETOPT_DUMP") == "1"
)

func main() {
	compmain.MaybeHandleCompletion()
	common.RunAndExit(realMain)
}

func print(args ...string) {
	for _, s := range args {
		os.Stdout.WriteString(s)
		if dump {
			os.Stderr.WriteString(s)
		}
	}
}

func printf(format string, args ...interface{}) {
	print(fmt.Sprintf(format, args...))
}

func realMain() int {
	printHeader()

	if err := getopt.CommandLine.Getopt(os.Args, nil); err != nil {
		common.Warn(err.Error())
		printExit(1)
		return 1
	}

	// Parse the flag spec.
	if len(getopt.Args()) == 0 {
		common.Warn("Missing option spec")
		printExit(1)
		return 1
	}

	specArg := getopt.Args()[0]
	spec, err := ParseOptionSpec(specArg)
	if err != nil {
		common.Warn(err.Error())
		printExit(1)
		return 1
	}

	if *bashCompletionCommand != "" {
		printRealBashcomp(spec, *bashCompletionCommand)
		return 0
	}

	// Build getopt options.
	options := getopt.New()
	for _, f := range spec {
		switch f.argType {
		case types.IsBoolean:
			options.FlagLong(&f.boolValue, f.longOption, f.shortOption, f.help)
		case types.IsInteger, types.IsString:
			options.FlagLong(&f.stringValue, f.longOption, f.shortOption, f.help)
		}
	}
	doHelp := false
	doBashComplete := false
	options.FlagLong(&doHelp, "help", 'h', "Show help")
	options.FlagLong(&doBashComplete, bashCompletionFlag, 0, "Print bash completion script")

	// Handle -h and parse errors.
	err = options.Getopt(getopt.Args(), nil)

	//common.Dump("Parsed=", spec)

	if err != nil {
		printEchoStderr(err.Error())
	}
	if err != nil || doHelp {
		printUsage(options)
		printExit(1)
		return 1
	}

	if doBashComplete {
		printBashcomp(specArg)
		printExit(0)
		return 0
	}

	// Validate arguments.
	for _, f := range spec {
		switch f.argType {
		case types.IsInteger:
			if f.stringValue != "" {
				if _, err = strconv.Atoi(f.stringValue); err != nil {
					printEchoStderr(fmt.Sprintf("\"%s\" is not integer", f.stringValue))
					printExit(1)
					return 1
				}
			}
		}
	}

	// Print results.
	for _, f := range spec {
		switch f.argType {
		case types.IsBoolean:
			if f.boolValue {
				//printf("# Short:%c  Long:%q\n", f.shortOption, f.longOption)
				print(f.evalString, "\n")
			}

		case types.IsInteger, types.IsString:
			if f.stringValue != "" {
				//printf("# Short:%c  Long:%q\n", f.shortOption, f.longOption)
				print(strings.Replace(f.evalString, "%", shell.Escape(f.stringValue), -1), "\n")
			}
		}
	}

	printNewArgs(options.Args())

	return 0
}

func printHeader() {
	print(`
# Succeeds if called by a function.
_go_infunc() {
  caller 1 >/dev/null 2>&1
}

# Get the command name, used for usage and completion.
_go_command=""
if _go_infunc ; then
  _go_command="$FUNCNAME"
else
  _go_command="${0##*\/}"
fi

`)
}

func printExit(rc int) {
	printf(`
# End with status %[1]d. 
if _go_infunc ; then
  return %[1]d
else
  exit %[1]d
fi

`, rc)
}

//func printEchoStdout(message string) {
//	printEchoInner(message, "")
//}

func printEchoStderr(message string) {
	printEchoInner(message, "1>&2")
}

func printEchoInner(message string, redirect string) {
	print("echo \"$_go_command: \"", shell.Escape(message), " ", redirect, "\n")
}

func printUsage(options *getopt.Set) {
	if *usage != "" {
		print(*usage, "\n")
	} else {
		printf(`
echo
echo "  $_go_command: "%s
echo
echo "  Usage:"
`, shell.Escape(*description))
	}

	b := &strings.Builder{}
	options.PrintOptions(b)
	print("echo ", shell.Escape(b.String()), "\n")
}

func printBashcomp(specArg string) {
	// Here, we want to show the script to install completion.
	// To do so, we need to know the command name, but we don't know that in the Go side.
	// We just know it in the Bash side as "$_go_command".
	// So, here, we print a command to invoke itself with
	// --bash-completion-command $_go_command.
	// When bashgetopt detects this flag, it then actually print the completion script.

	dashN := ""
	if *noAllowFiles {
		dashN = "-N"
	}
	printf(
		`%s %s --`+realBashCompletionCommandFlag+` "$_go_command" %s`,
		shell.Escape(common.MustGetExecutable()),
		dashN,
		shell.Escape(specArg))
}

func printRealBashcomp(spec []*OptionSpec, command string) {
	// Create completion spec.
	specStr := bytes.Buffer{}
	specStr.WriteString("@switchloop \"^-\"\n")
	for _, f := range spec {
		if f.longOption == bashCompletionFlag {
			continue
		}

		addArg := func() {
			switch f.argType {
			case types.IsInteger:
				specStr.WriteString(fmt.Sprintf("      @any # [INTEGER] %s\n", f.help))
			case types.IsString:
				specStr.WriteString(fmt.Sprintf("      @any # [STRING] %s\n", f.help))
			}
		}

		if f.shortOption != 0 {
			specStr.WriteString(fmt.Sprintf("    -%c # %s\n", f.shortOption, f.help))
			addArg()
		}
		if f.longOption != "" {
			specStr.WriteString(fmt.Sprintf("    --%s # %s\n", f.longOption, f.help))
			addArg()
		}
	}
	if !*noAllowFiles {
		specStr.WriteString("@loop\n")
		specStr.WriteString("    @cand takeFile\n")
	}

	all := bytes.Buffer{}

	opts := compmain.InstallOptions{In: os.Stdin, Out: &all}

	compmain.PrintInstallScriptRaw(specStr.String(), opts, command)

	print(all.String())
}

func printNewArgs(args []string) {
	print("set -- ")
	for i, arg := range args {
		if i > 0 {
			print(" ")
		}
		print(shell.Escape(arg))
	}
	print("\n")
}
