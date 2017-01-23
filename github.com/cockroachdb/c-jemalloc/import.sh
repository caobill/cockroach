#!/usr/bin/env sh

MAKE=${MAKE:-make}
set -eu

rm -rf internal/*
find . -type l -not -path './.git/*' -exec rm {} \;
curl -sL https://github.com/jemalloc/jemalloc/releases/download/4.4.0/jemalloc-4.4.0.tar.bz2 | tar jxf - -C internal --strip-components=1
patch -p1 -d internal < secure_getenv.patch

# You need to manually run the following code.
# on OSX:
# (cd internal && MACOSX_DEPLOYMENT_TARGET=10.9 ./configure --enable-prof --with-jemalloc-prefix='')
# <compare "Build parameters" in internal/Makefile to cgo flags in cgo_flags.go> and adjust the latter.
# rm -r darwin_includes
# git clean -Xn -- internal/include/jemalloc | sed 's/.* //' | xargs -I % rsync -R % darwin_includes/
#
# on Linux:
# cd internal
# echo 'ac_cv_func_issetugid=no' >> config.cache
# echo 'ac_cv_func_secure_getenv=no' >> config.cache
# echo 'je_cv_glibc_malloc_hook=no' >> config.cache
# echo 'je_cv_glibc_memalign_hook=no' >> config.cache
# echo 'je_cv_madv_free=no' >> config.cache
# echo 'je_cv_pthread_mutex_adaptive_np=no' >> config.cache
# echo 'je_cv_thp=no' >> config.cache
# ./configure --enable-prof --disable-prof-libgcc -C
# rm config.cache
# cd -
# <compare "Build parameters" in internal/Makefile to cgo flags in cgo_flags.go> and adjust the latter.
# rm -r linux_includes
# git clean -Xn -- internal/include/jemalloc | sed 's/.* //' | xargs -I % rsync -R % linux_includes/
#
# on FreeBSD:
# (cd internal && ./configure --enable-prof)
# <compare "Build parameters" in internal/Makefile to cgo flags in cgo_flags.go> and adjust the latter.
# rm -r freebsd_includes
# git clean -Xn -- internal/include/jemalloc | sed 's/.* //' | xargs -I % rsync -R % freebsd_includes/

# symlink so cgo compiles them
for source_file in $($MAKE sources); do
  ln -sf "$source_file" .
done

# restore the repo to what it would look like when first cloned.
# comment this line out while updating upstream.
git clean -dxf
