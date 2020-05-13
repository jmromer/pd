pd
==

Fast, fuzzy switching between version-controlled projects and directories on the
command line.

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
    local target

    if [[ "$1" =~ ^[-+][0-9]+ ]]; then
        target="$1"
    elif [[ -z "$1" ]]; then
        target="$(pd)"
    else
        target="$(pd cd "$1")"
    fi

    [[ -z "$target" ]] && return
    builtin cd "$target" || return
}

# ctrl-o to change directories with pd
bindkey -s '^o' 'cd\n'
```

[ascii-svg]: https://asciinema.org/a/sqrGsf4drptaOyU6UUaJ4OSgN.svg
[ascii]: https://asciinema.org/a/sqrGsf4drptaOyU6UUaJ4OSgN
