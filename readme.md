# Usage:

- Visit `zsync http://metadata.neon.kde.org/snap/plasma_5.9_amd64.snap.zsync` to acquire the snap package.
- Execute `sudo snap install --force-dangerous --devmode plasma_5.9_amd64.snap` to install the snap package.
- This is only necessary when the snap was newly installed/reinstalled:
  - Execute `snap list` to observe any installed revisions of `plasma-desktop` (modifications during each installation) such as `x8`.
  - Copy the `bin/` directory from the repository into `$HOME/snap/plasma/REVISION/`.
- Execute (currently runs a bash) `dbus-launch snap run plasma`
- Configure environment `$HOME/bin/envo`.
- Now, startplasmacompositor, etcetera, should be operational.
- One simple test in `xsession` `kwin_wayland --xwayland` should create one black Wayland window.

# Miscellany:

- The `bin/` stuff is a hack so that the snap is not requiring to be reconstruction when the environment is requiring expansion, etcetera.
- This is using one patched `xwayland`. `xkb/ddxLoad.c` in `xorg` is hardcompiling the `xkbcomp` path to `/usr/bin`. The patched build simply strips the pathing so that `xkbcomp` is resolved from `$PATH` (this is not secure). There is a relocatibility feature resolving the location of `xkbcomp` relative to the server binary, but it is only implemented and being used by Windows Server currently.
- `kinit` hard-compile paths in `start_kdeinit` and `start_kdeinit_wrapper`. The entire pile of binaries is appearing to be somewhat weird because in part is claims `setuid` is involved yet none of the involved binaries actually has the `setuid` set apparently. `start_kdeinit` does however `oom_adj`, whether or not that actually is operational is unknown, because it has no `setuid` *and* it is appearing to utterly ignore return values, so I am guessing that it is not able to operate, probably.
- `kinit` is having patches to turn the hard-compiled paths into lookup. For wrapper this lookup is relative within `same-dir`, whereas for `start_kdeinit` (forks `kdeinit`) is requiring redirection through `$PATH`. This would be a hazardous securitywise if it was `setuid` and passed that own to any forks, but because that it is not to begin with, this is probably acceptable.
