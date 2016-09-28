Go support for EDF+
===================

This package attempt to provide a Go implementation of the EDF format. It reads EDF+ files into two structures:

+ A map from strings to strings, representing the main header in the EDF file.
+ A slice of slices of int16, each one representing one channel's recording.

They can be accessed using the `ReadFile(string)` function. 

To Do
-----

- [x] Organize folder in Go style.
- [x] Check how to import C code from Go.
- [ ] Write documentation in Go style.