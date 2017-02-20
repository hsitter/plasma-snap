// -*- Mode: Go; indent-tabs-mode: t -*-

/*
	Copyright 2017 Harald Sitter <sitter@kde.org>

	This program is free software; you can redistribute it and/or
	modify it under the terms of the GNU General Public License as
	published by the Free Software Foundation; either version 3 of
	the License or any later version accepted by the membership of
	KDE e.V. (or its successor approved by the membership of KDE
	e.V.), which shall act as a proxy defined in Section 14 of
	version 3 of the license.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package builtin

import (
	"bytes"
	"fmt"

	"github.com/snapcore/snapd/interfaces"
)

const plasmaConnectedSlotAppArmor = `
`

const plasmaPermanentSlotDBus = `
`

// FIXME: we are basically the wrong way around here. the slot should have these
// abilities, the plug shouldn't hue hue
const plasmaConnectedPlugAppArmor = `
# Description: Magic

# TODO: why?
/sys/devices/*/*/boot_vga r,

# Possibly restrict this more. KWin diables tracing on itself.
ptrace (trace),

# xwayland socket
# bind cannot be used with peer QQ
unix (connect, receive, send, bind, listen)
     type=stream,
#     peer=(addr="@/tmp/.X11-unix/*"),

# kcrash checks pattern
/proc/sys/kernel/core_pattern r,

# krunner iterates this for some reason
# TODO: figure out why
/etc/fstab r,
/sys/bus/*/devices/{,*} r,
/run/udev/data/{,*} r,

#include <abstractions/dbus-strict>

# Used in various places to check if a thing is registered before calling it.
# Most importantly used in kwin before beginning to talk to login1
dbus (send)
    bus={system,session}
    path={/,/org/freedesktop/DBus}
    interface=org.freedesktop.DBus
    member=ListNames
    peer=(name="org.freedesktop.DBus", label="unconfined"),

# KWin manages session activity and uses login1 to use devices.
dbus (send)
    bus=system
    path=/org/freedesktop/login1{,/**}
    interface=org.freedesktop.login1.Manager
    member={GetSession,GetSessionByPID}
    peer=(name=org.freedesktop.login1),
## Next rule can't be locked to peer becuase PauseDevice is weirdly shot against :1.0
dbus (send, receive)
    bus=system
    path=/org/freedesktop/login1/session/*
    interface=org.freedesktop.login1.Session
    member={Activate,TakeControl,TakeDevice,ReleaseDevice,PauseDevice},
dbus (send)
    bus=system
    path=/org/freedesktop/login1/session/*
    interface=org.freedesktop.DBus.Properties
    member=Get{,All}
    peer=(name=org.freedesktop.login1),

# Powerdevil wants to inhibit
dbus (send)
    bus=system
    path=/org/freedesktop/login1
    interface=org.freedesktop.login1.Manager
    member=Inhibit
    peer=(name=org.freedesktop.login1),

# Introspection through solid. This is used in plugable tech.
# TODO: why does the udisks2 plug not allow this?
dbus (send)
    bus=system
    path=/
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(name=org.freedesktop.UDisks2),
dbus (send)
    bus=system
    path=/
    interface=org.freedesktop.DBus.Properties
    member=Get{,All}
    peer=(name=org.freedesktop.UDisks2),
dbus (send)
    bus=system
    path=/org/freedesktop/UDisks2/{block_devices,drives}
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(name=org.freedesktop.UDisks2),
dbus (send)
    bus=system
    path=/org/freedesktop/UDisks2/{block_devices,drives}/*
    interface=org.freedesktop.DBus.Properties
    member=Get{,All}
    peer=(name=org.freedesktop.UDisks2),
dbus (send)
    bus=system
    path=/org/freedesktop/UDisks2/{block_devices,drives}/*
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(name=org.freedesktop.UDisks2),
# TODO: upower plug also seems a bit meh
dbus (send)
    bus=system
    path=/org/freedesktop/UPower
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(name=org.freedesktop.UPower,label=unconfined),

# plasma-nm (&kded) allow full access more or less.
dbus (send)
    bus=system
    path=/org/freedesktop/NetworkManager{,/**}
    interface=org.freedesktop.DBus.Properties
    member=Get{,All}
    peer=(name=org.freedesktop.NetworkManager),
dbus (send)
    bus=system
    path=/org/freedesktop/NetworkManager{,/**}
    interface=org.freedesktop.NetworkManager{,.**}
    peer=(name=org.freedesktop.NetworkManager),
`

// PlasmaInterface is the hello interface for a tutorial.
type PlasmaInterface struct{}

// String returns the same value as Name().
func (iface *PlasmaInterface) Name() string {
	return "plasma"
}

// SanitizeSlot checks and possibly modifies a slot.
func (iface *PlasmaInterface) SanitizeSlot(slot *interfaces.Slot) error {
	if iface.Name() != slot.Interface {
		panic(fmt.Sprintf("slot is not of interface %q", iface))
	}
	// NOTE: currently we don't check anything on the slot side.
	return nil
}

// SanitizePlug checks and possibly modifies a plug.
func (iface *PlasmaInterface) SanitizePlug(plug *interfaces.Plug) error {
	if iface.Name() != plug.Interface {
		panic(fmt.Sprintf("plug is not of interface %q", iface))
	}
	// NOTE: currently we don't check anything on the plug side.
	return nil
}

// ConnectedSlotSnippet returns security snippet specific to a given connection between the hello slot and some plug.
func (iface *PlasmaInterface) ConnectedSlotSnippet(plug *interfaces.Plug, slot *interfaces.Slot, securitySystem interfaces.SecuritySystem) ([]byte, error) {
	switch securitySystem {
	case interfaces.SecurityAppArmor:
		old := []byte("###SLOT_SECURITY_TAGS###")
		new := slotAppLabelExpr(slot)
		snippet := bytes.Replace([]byte(plasmaConnectedSlotAppArmor), old, new, -1)
		return snippet, nil
	}
	return nil, nil
}

// PermanentSlotSnippet returns security snippet permanently granted to hello slots.
func (iface *PlasmaInterface) PermanentSlotSnippet(slot *interfaces.Slot, securitySystem interfaces.SecuritySystem) ([]byte, error) {
	switch securitySystem {
	case interfaces.SecurityDBus:
		return []byte(plasmaPermanentSlotDBus), nil
	}
	return nil, nil
}

// ConnectedPlugSnippet returns security snippet specific to a given connection between the hello plug and some slot.
func (iface *PlasmaInterface) ConnectedPlugSnippet(plug *interfaces.Plug, slot *interfaces.Slot, securitySystem interfaces.SecuritySystem) ([]byte, error) {
	switch securitySystem {
	case interfaces.SecurityAppArmor:
		old := []byte("###SLOT_SECURITY_TAGS###")
		new := slotAppLabelExpr(slot)
		snippet := bytes.Replace([]byte(plasmaConnectedPlugAppArmor), old, new, -1)
		return snippet, nil
	}
	return nil, nil
}

// PermanentPlugSnippet returns the configuration snippet required to use a hello interface.
func (iface *PlasmaInterface) PermanentPlugSnippet(plug *interfaces.Plug, securitySystem interfaces.SecuritySystem) ([]byte, error) {
	switch securitySystem {
	case interfaces.SecurityAppArmor:
		return nil, nil
	case interfaces.SecuritySecComp:
		return []byte(`bind`), nil
	}
	return nil, nil
}

// AutoConnect returns true if plugs and slots should be implicitly
// auto-connected when an unambiguous connection candidate is available.
//
// This interface does not auto-connect.
func (iface *PlasmaInterface) AutoConnect(*interfaces.Plug, *interfaces.Slot) bool {
	return false
}

func init() {
	allInterfaces = append(allInterfaces, &PlasmaInterface{})
}
