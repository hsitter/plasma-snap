# Use

- get snap with `zsync http://metadata.neon.kde.org/snap/plasma_5.9_amd64.snap.zsync`
- install `sudo snap install --force-dangerous --devmode plasma_5.9_amd64.snap`
- this is only necessary when the snap was newly installed/reinstalled
  - check `snap list` for installed rev of plasma snap (changes on each install). e.g. `x8`
  - copy the `bin/` directory from the repo into $HOME/snap/plasma/REVISION/
- run (currently runs a bash) `dbus-launch snap run plasma`
- setup environment `$HOME/bin/envo`
- at this point startplasmacompositor would work etc.
- simple test in xsession `kwin_wayland --xwayland` should bring up a black wayland window

# Misc

- the bin/ stuff is a hack so the snap doesn't have to rebuilt every time the environment needs to be expanded etc.
- this uses a patched xwayland. xkb/ddxLoad.c in xorg is hardcompiling the xkbcomp path to /usr/bin, the patched build simply strips the pathing so xkbcomp is resolved from $PATH (nb: this is a security leak). There is a relocatibility feature resolving xkbcomp's location relative to the server binary, but it is only implemented and used in the windows server right now.
