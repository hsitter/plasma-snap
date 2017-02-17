#!/usr/bin/env ruby

require '/var/lib/jenkins/ci-tooling/nci/setup_apt'

ENV['TYPE'] = 'release'

NCI.setup_repo!

system('apt update') || raise
system('apt install -y snapcraft git zsync') || raise
