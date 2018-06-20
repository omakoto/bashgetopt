[![Build Status](https://travis-ci.org/omakoto/bashgetopt.svg?branch=master)](https://travis-ci.org/omakoto/bashgetopt)
# Bashgetopt

A command line option parser for Bash scripts. 

## Usage

See [the sample](sample/pg).

Essentially, the definition would look like this:
```bash
eval "$(bashgetopt -d "Poor grep" '
i|ignroe-case    ignorecase=1    # ignore case distinctions
C|conext=i       context=%       # print NUM lines of context
e|regexp=s       pattern=%       # use PATTERN for matching
q                quiet=1         # don'\''t print matching lines. (<- this is how to put single quotes)
' "$@")"
```
 - The first column defines flags and types
    (`=s` for flags with string arguments and `=i` for integer arguments).
 - Everything after `#` is used for help.
 - The middle part will is what'll be executed when each flag is detected.
   An argument will be replaced with `%`.

## Showing help

`bashgetopt` automatically will generate help if `-h` or `--help` is given.

```bashs   
$ ./sample/pg -h

  pg: Poor grep

  Usage:
     --bash-completion
                    Print bash completion script
 -C, --conext=value
                    print NUM lines of context
 -e, --regexp=value
                    use PATTERN for matching
 -h, --help         Show help
 -i, --ignroe-case  ignore case distinctions
 -q                 don't print matching lines. (<- this is how to put single
                    quotes)
```

## Shell completion

`bashgetopt` supports shell completion with [compromise](https://github.com/omakoto/compromise)
for Bash and Zsh. Add the following line to your `.bashrc` or `.zshrc`. 

```bash
. <(pg --bash-completion)
```

## Installation

go get -u github.com/omakoto/bashgetopt
