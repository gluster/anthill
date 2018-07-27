This page provides an overview of the project's testing infrastructure.

# Git pre-commit hook

The scripts directory has a script that can be used as a git pre-commit hook to
lint your files prior to committing. While this is totally optional, our CI
infrasturcture will run it and fail your PR if it doesn't pass cleanly.

You can run the script manually:

```sh
$ ./scripts/pre-commit.sh

mdl ./.github/ISSUE_TEMPLATE/bug_report.md
mdl ./.github/ISSUE_TEMPLATE/feature_request.md
mdl ./.github/pull_request_template.md
mdl ./docs/index.md
mdl ./docs/Developers/Testing_and_validation.md
mdl ./docs/Developers/Editing_the_documentation.md
mdl ./docs/Users_guide/index.md
mdl ./CONTRIBUTING.md
mdl ./README.md
shellcheck ./scripts/pre-commit.sh
yamllint -s ./.travis.yml
yamllint -s ./mkdocs.yml
ALL OK.
```

Optionally, you can install it to run automatically each time you `git commit`.
To install the hook, create a link to the script in your local `.git/hooks`
directory:

```sh
cd .git/hooks
ln -s ../../scripts/pre-commit.sh pre-commit
```

# CI infrastructure

The project makes use of Travis CI for linting (and soon, unit testing).

A check by Travis (see `.travis.yml`) is invoked with each commit and PR.
Currently, it:

- runs the pre-commit.sh script to lint all formatted files
- runs mkdocs to ensure the documentation builds cleanly
