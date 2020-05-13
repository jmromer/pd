pd
==

Fast, fuzzy switching between projects on the command line.

[![asciicast][ascii-svg]][ascii]

Recommended setup
-----------------

```sh
# ~/.zshrc

# override built-in cd to:
# 1. fuzzy-select a directory to visit (when given no args)
# 2. retain built-in dir-stack-related behavior (when given a leading -/+ numeric arg)
# 3. log a directory visit (when given any other arg)
cd() {
    local target
    case "$1" in
        "")
            target="$(pd)"
            ;;
        ^[-+][0-9]+$)
            target="$1"
            ;;
        *)
            target="$(pd cd "$1")"
            ;;
    esac

    [[ -z "$target" ]] && return
    builtin cd "$target" || return
}

# ctrl-o to change directories with pd
bindkey -s '^o' 'cd\n'
```

[ascii-svg]: https://asciinema.org/a/sqrGsf4drptaOyU6UUaJ4OSgN.svg
[ascii]: https://asciinema.org/a/sqrGsf4drptaOyU6UUaJ4OSgN
