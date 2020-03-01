# `e`: An environment variable manager

`e` Allows you to define environment profiles and load/switch environment variables easily

Here's how the `e` journey starts:

```
e create perso        # Create a personal profile
e perso               # Select the perso profile
e set CLIENT_ID 1234  # Set an env var
echo $CLIENT_ID       # The value is propagated to your current shell

e create pro          # Create a new profile, for profesional things
e pro                 # Activate the pro profile
echo $CLIENT_ID       # This will print nothing, the value was unset!
e set CLIENT_ID 5678  # Set the same env var, in the new profile
```

The profiles are simple `.env` type files saved under `~/.e`


## How does it work

The `e` binary only takes care of the profile managment. It creates, lists, updates the files under `~/.e`. When a command changes the selected profiles or any of its values, the binary outputs the new profile environment variables. A shell wrapper is used to read those variables and propagate the values to the current shell.

For example, the PowerShell wrapper reads all the variables returned by the binary and calls `[Environment]::SetEnvironmentVariable` with the correct key and value