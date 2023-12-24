# bulk-encrypted-zip-extractor

A small script that uses 7Z binary to extract multiple zip/rar/7z files getting passwords from a password.txt file

# Requirements

- `7zz` binary from https://www.7-zip.org/download.html

# Usage:

- `./bulk-encrypted-zip-extractor <dir_with_zip_files> <path_to_7zz> <path_to_passwords.txt> <true/false (delete zip after extracting>`

# Compile:

- `go build -trimpath .`