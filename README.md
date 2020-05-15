p/d
===

A project / directory manager and FZF-powered fuzzy-selector.

Issue `pd` to search your home directory for version-controlled projects and
index them. You'll be able to fuzzy-select from this list, which is ranked in
order of visit frequency. Use `pd` in tandem with `cd` (see below) to add new
directories to the index, and log new visits.


[![asciicast][ascii-svg]][ascii]

Recommended setup
-----------------

```sh
# ~/.zshrc

# wrap built-in cd to:
# 1. fuzzy-select a directory to visit when given no arg
# 2. retain built-in dir-stack-related behavior when given a leading -/+ numeric arg
# 3. log a directory visit when given any other arg

cd() {
    builtin cd "$(pd "$1")" || return
}

# ctrl-o to change directories with pd
bindkey -s '^o' 'cd\n'
```

[ascii-svg]: https://asciinema.org/a/330578.svg
[ascii]: https://asciinema.org/a/330578
