#!/usr/bin/env bash

archi=`arch`
if [ "$archi" == "x86_64" ];then
  archi="amd64"
elif [ "$archi" == "i386" ];then
  archi="arm64"
fi

os=`uname | tr '[A-Z]' '[a-z]'`

package=`curl -s https://api.github.com/repos/vine-io/gpm/releases/latest | grep browser_download_url | grep ${os} | cut -d'"' -f4 | grep "gpm-${os}-${archi}"`

echo "install package: ${package}"
wget ${package} -O /tmp/gpm.tar.gz && mkdir -pv /tmp/gpm && tar -xvf /tmp/gpm.tar.gz -C /tmp/gpm

mv /tmp/gpm/${os}/* /usr/sbin/

rm -fr /tmp/gpm
rm -fr /tmp/gpm.tar.gz
