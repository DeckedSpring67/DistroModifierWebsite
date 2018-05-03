#!/bin/bash
ARCH_DIR=../arch.temp/arch/x86_64
umount -lRf $ARCH_DIR/base-root
umount -lRf $ARCH_DIR/base-root/*
mksquashfs $ARCH_DIR/base-root $ARCH_DIR/airootfs.sfs
rm -rf $ARCH_DIR/base-root
gpg --output $ARCH_DIR/airootfs.sfs.sig --sign $ARCH_DIR/airootfs.sfs
sha512sum $ARCH_DIR/airootfs.sfs > $ARCH_DIR/airootfs.sha512
createiso.sh $1
