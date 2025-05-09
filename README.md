# pre-commit-check-coverage

A git pre-commit hook to ensure that staged files are covered by tests before allowing commits.

It requires a textfile sha256sums.txt in the repository root formatted like `(<sha256sum>\t<filename>)*` with files included in the coverage report.

It will then match all files in the changeset with this list to identify any files where the hash and therefore content does not match (the implicit assumption being that the current version of the file has not been tested).