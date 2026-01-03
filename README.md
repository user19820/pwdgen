# pwdgen

Simple password manager. Generates random passwords of specified length (default 20 chars)
and stores them with a given name. The passwords can then be retrieved later using the same name and
stored in your system clipboard.

pwdgen should support darwin, windows and linux, though linux and windows support is currently untested.

## Commands

- `init`: sets up pwdgen
- `gen <name> [<length>]`: generates a random password, which is stored with _name_ 
and has _length_ number of chars. _length_ is optional and defaults to 25 characters.
- `get <name>`: retrieves the password with _name_, if it exists.

