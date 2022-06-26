#!/bin/bash
# User defined commands

deploy_backoffice() {
    cd /code || return
    make GIT_TAG="dev" build/linux/backoffice || return 1
    raftt stop backoffice
    raftt cp build/linux/backoffice backoffice:/app/build/linux/backoffice
    raftt restart backoffice
}