# Pwdgen

Simple password generator. Uses an extended character set and imposes minimum
length requirement to (hopefully) achieve decently strong passwords.

IMPORTANT NOTE: Results are directly copied to the user's clipboard.

Supports:
- darwin
- windows (needs testing!)
- linux (needs testing!)

TODO:
- research better password generation algorithms

Note on implementation:
In this project I am experimenting with the use of OS-specific compilation.
Basically, `copyToClipboard` is implemented separately for macOS, Windows and linux
and during comptime the compiler picks the implementation which allows the program to
run.

