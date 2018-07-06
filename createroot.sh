#!/bin/bash
WORKDIR=$(pwd)
echo $WORKDIR
cd ..
rm -rf arch.temp
cp -rf arch arch.temp
rm -rf arch.temp/arch/x86_64/*
cp -rf base-root arch.temp/arch/x86_64/
cd arch.temp/arch/x86_64/base-root/
mount -t proc /proc proc/
mount --rbind /dev/ dev/
mount --rbind /sys/ sys/
mount --rbind . .
chroot . pacman -Syy --noconfirm
for x in $@
do
if (( $x == gnome ))
then
	chroot . systemctl enable gdm
fi
chroot . pacman --noconfirm -S $x 
done
chroot . pacman -U linux-4.15.14-1-x86_64.pkg.tar.xz --noconfirm
chroot . pacman -Scc
chroot . chown -R live /home/live/

