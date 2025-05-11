# pre-commit-check-coverage

A git pre-commit hook that blocks committing untested code.

This is done using a combination of a test coverage report and the source files.

It will then, based on file name and content, check and fail if a staged file is not
present in the coverage report as this implies the test suite has not been run with the current version of the source file.

## Custom Setup/Pre-requisites

It requires a textfile sha256sums.txt in the repository root formatted like `(<sha256sum>\t<filename>)*` with files included in the coverage report. 

This file must be re-computed/refreshed each time a new report is generated (or the check is performed)
as it is not done by the check.

This is done rather easily by wrapping the checksum producing code around the report generator code.

## How it works

### In Theory

Comparing checksum of files is a long standing method to check for sameness.

For each staged file, its checksum is computed.

Similarly, checksums are computed for all files included in the coverage report.

The check fails if `(checksum staged file x filename)\(checksum tested file x filename) = Ø` does not hold.


### In Practice

Note: Eventually it is possible to write parsers or a plugin system to read coverage reports but as a start, the exercise is left for the user to generate the logic that computes checksums for a report. 

Checksum and filename is required information, as the checksum alone does not properly allow verification as different files can have the same content and therefore the same checksum.

Ideally, the checksum file is updated every time the test suite runs but it can be computed just in time of running the check to supply a valid response.

`git diff-index` provides a list of file changes along with their status. Added, modified and renamed files are included - deleted files are skipped.

The checksum of the file in the worktree is then computed and matched with the checksums and filenames in the coverage report.


## Future work/features

- Checking against unstaged code changes as this can yield false negatives
- Separate cli interface from module interface
- Pattern based match/filtering - not all staged files are code changes
- pre-commit hook definition
- github actions running code example