Go support for EDF+
===================

This package attempt to provide a Go implementation of the EDF format. It reads EDF+ files into two structures:

+ A map from strings to strings, representing the main header in the EDF file.
+ A slice of slices of int16, each one representing one channel's recording.
