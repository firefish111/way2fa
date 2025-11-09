# The `.way` format
This is the custom format that `way2fa` uses to store 2FA keys. This is its specification.

### Pure formats
`way2fa` does, as a historical accident, support some formats that do not adhere to this specification: these are referred to in the codebase as "pure" formats, and they ought to be disregarded as they are awfully insecure.

## Format specification

The format consists of the following three consecutive fields, each seperated by a `.` character.

| Field Name | Size | Purpose |
| :- | :-: | :- |
| Header | `8` bytes | To signify that this is indeed a `.way` file, and which module to use |
| Encrypted header | `32` bytes | To verify that the inserted password was correct |
| Content | ? | Whatever content the file is housing |

# To be continued
