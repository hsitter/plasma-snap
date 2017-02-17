#!/usr/bin/env ruby

require '/var/lib/jenkins/ci-tooling/nci/setup_apt'

system('apt update') || raise
system('apt install -y snapcraft git zsync') || raise
