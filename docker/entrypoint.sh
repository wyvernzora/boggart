#!/bin/sh

if [ ! -f /boggart/etc/host_key ]; then
  echo '/boggart/etc/host_key not found, generating'
  ssh-keygen -t ed25519 -N '' -f /boggart/etc/host_key
end

/bin/boggart /boggart/etc/config.yml
