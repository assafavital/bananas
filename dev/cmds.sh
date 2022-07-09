#!/bin/bash
# User defined commands

deploy_backoffice() {
    cd /code || return
    make GIT_TAG="dev" build/linux/backoffice
    raftt stop backoffice
    raftt sh backoffice -- sh -c 'cat /dev/null > /tmp/log/raftt/lifeguard/hijacked.out'
    raftt cp build/linux/backoffice backoffice:/code/build/linux/backoffice
    raftt restart backoffice
    raftt logs backoffice -f | /root/.cargo/bin/fblog
}