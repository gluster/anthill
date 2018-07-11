#! /bin/bash
# vim: set ts=4 sw=4 et :

set -e

# Run checks from root of the repo
scriptdir="$(dirname "$(realpath "$0")")"
cd "$scriptdir/.."

# run_check <file_regex> <checker_exe> [optional args to checker...]
function run_check() {
    regex="$1"
    shift
    exe="$1"
    shift

    if [ -x "$(command -v "$exe")" ]; then
        find . -regextype egrep -iregex "$regex" -print0 | \
            xargs -0rt -n1 "$exe" "$@"
    else
        echo "Warning: $exe not found... skipping some tests."
    fi
}


# Install via: gem install asciidoctor
run_check '.*\.adoc' asciidoctor -o /dev/null -v --failure-level WARN

# markdownlint: https://github.com/markdownlint/markdownlint
# https://github.com/markdownlint/markdownlint/blob/master/docs/RULES.md
# Install via: gem install mdl
run_check '.*\.md' mdl

# Install via: dnf install shellcheck
run_check '.*\.(ba)?sh' shellcheck

# Install via: pip install yamllint
run_check '.*\.ya?ml' yamllint -s

echo "ALL OK."
