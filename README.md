# govault

A super simple commandline based encrypted vault to store secret data with a key for easy copying to clipboard

## Usage

```
:h -> help
:q -> quit
:c -> clear screen
:k -> show keys
:s -> show key/value
:d key -> delete key
key<SPACE>value -> create/update entry
key<ENTER> -> copy value to clipboard
```

## Build

```
make all
```

## Run

```
release/govault-<os>-<version>
Enter master password:
```

A vault file will be created in $HOME/.govault by default.
You can change location of vault file by setting `GO_VAULT_FILE` env.

To add/update a key/value to vault, type `key<SPACE>value` and it will be added to vault file encrypted with your master password.

To copy value of a key to clipboard, type `key<ENTER>`.
