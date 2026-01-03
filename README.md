# pwdgen

Simple password manager. Generates random passwords of specified length (default 20 chars)
and stores them with a given name. The passwords can then be retrieved later using the same name and
stored in your system clipboard.

pwdgen should support darwin, windows and linux, though linux and windows support is currently untested.

## Commands

- `init`: sets up pwdgen
- `gen <name> [<length>]: generates a random password, which is stored with name <name>
and has <length> number of chars. <length> is optional and defaults to 25 characters.
- `get <name>: retrieves the password with <name>, if it exists.

