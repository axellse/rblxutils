@echo off
echo "updating rblxutils... dont close this window"
timeout /t 3
del rblxutils.exe
ren rblxutils_fresh.exe rblxutils.exe
start rblxutils.exe