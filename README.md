p/d
===

A project / directory manager and [FZF][]-powered fuzzy-selector.

Use `pd` in tandem with `cd` to change directories using fuzzy-search, search
for version-controlled projects in your home directory, and keep track of your
most-visited `cd` targets.

It's written in [Go][], and it _zooms_.

[![asciicast][ascii-svg]][ascii]

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Contents**

- [Recommended setup](#recommended-setup)
- [Usage](#usage)
- [Installation](#installation)
- [License](#license)
- [Acknowledgements](#acknowledgements)

<!-- markdown-toc end -->

Recommended setup
-----------------

```sh
# ~/.zshrc

# wrap built-in cd to:
# 1. fuzzy-select a directory to visit when given no argument
# 2. retain built-in dir-stack-related behavior when given a -/+ numeric arg
# 3. log a directory visit when given any other arg

cd() {
    builtin cd "$(pd "$1")" || return
}

# ^h cd with pd
bindkey -s '^h' 'cd\n'
```

Usage
-----
```
% pd --help

p/d

Usage:

  pd [directory name]

Intended to be used in tandem with cd as follows:

  cd $(pd ~/Documents)

Given a file path, print its absolute form (resolving symlinks) and save to
history. If a path to a non-directory is given, use its containing
directory instead.

Examples:

  pd ~/Documents/projects/my-project
  pd ~/my-other-project
  pd ./projects/my-project/some-file.txt

Given a position on the directory stack, no-op. Print that back out to leave the
behavior of cd unchanged.

Examples:
  pd -2
  pd +1
  pd -

Given no arguments, open FZF to allow fuzzy-selecting a directory to cd into.
```

Installation
------------

Still a work in progress, so not yet published.

Clone and build with `go build && go install`.

License
-------

Apache

Acknowledgements
----------------

p/d is written in [Go][] based on a prototype in [Ruby][].
It builds upon prior art by [junegunn][] ([fzf][]) and [b4b4r07][]
([go-finder][]).

[ascii-svg]: https://asciinema.org/a/330647.svg
[ascii]: https://asciinema.org/a/330647
[b4b4r07]: https://github.com/b4b4r07
[fzf]: https://github.com/junegunn/fzf
[go-finder]: https://github.com/b4b4r07/go-finder
[go]: https://golang.org/doc
[junegunn]: https://github.com/junegunn
[ruby]: https://ruby-doc.org/
