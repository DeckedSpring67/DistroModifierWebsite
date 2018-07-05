#!/bin/sh
cd ../arch.temp
CDIR=$(pwd)
echo "Current Dir= $CDIR"
rm -rf arch/x86_64/base/base-root
xorriso -as mkisofs \
  -iso-level 3 \
  -isohybrid-mbr "isolinux/isohdpfx.bin" \
  -c "isolinux/boot.cat" \
  -b "isolinux/isolinux.bin" \
  -A "ARCH_201804" \
  -volid "ARCH_201804" \
  -no-emul-boot \
  -boot-load-size 4 \
  -boot-info-table \
  -eltorito-alt-boot \
  -e "EFI/archiso/efiboot.img" \
  -no-emul-boot \
  -isohybrid-gpt-basdat \
  -o ../$1.iso \
   "$CDIR"
cd ..
mv $1.iso /var/www/localhost/htdocs/
echo "done"


