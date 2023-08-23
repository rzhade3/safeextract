# Safe Extract

## Description

This is a simple library to safely extract files from a `zip` or `tar` archive. It checks:

* Archive only contains regular files and directories (or safe symlinks if so configured)
* No files are extracted outside of the target directory
* No files are extracted with absolute paths
* No symlinks are extracted which point outside of the target directory
* The archive is smaller than a configurable limit
