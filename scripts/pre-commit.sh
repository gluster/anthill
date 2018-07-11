#! /bin/bash
# vim: set ts=4 sw=4 et :

set -e

# Run checks from root of the repo
scriptdir="$(dirname "$(realpath "$0")")"
cd "$scriptdir/.."

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

run_check '.*\.adoc' asciidoctor -o /dev/null -v --failure-level WARN
run_check '.*\.md' mdl
run_check '.*\.(ba)?sh' shellcheck
run_check '.*\.ya?ml' yamllint -s

echo "ALL OK."
