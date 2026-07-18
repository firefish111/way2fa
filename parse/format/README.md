# The `.way` format
The `.way` format is the custom format that `way2fa` uses to store 2FA keys. This is its specification.

### Pure formats
`way2fa` does, as a historical accident, support some formats that do not adhere to this specification, and instead store data as plain text; these are referred to in the codebase as "pure" formats, and they ought to be disregarded as they are awfully insecure. The UI will warn you on startup if you were to use a pure format.

## Format specification

The format consists of the following consecutive fields.

| Field Name | Size | Purpose |
| :- | :-: | :- |
| Header | `8` bytes | To signify that this is indeed a `.way` file, and which module to use |
| IV | `12` bytes | Cryptographic nonce for AES encryption. All nulls if no password. |
| Salt | `16` bytes | Cryptographic salt for key derivation. All nulls if no password. |
| Content | ? | Whatever content the file is housing, base64encoded then encrypted in algorithm signified in header. |

# To be continued
