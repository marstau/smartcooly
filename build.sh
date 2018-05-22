#!/bin/bash

xgo --targets=windows/*,darwin/amd64,linux/amd64,linux/386,linux/arm --dest=cache ./

osarchs=(windows_amd64 windows_386 darwin_amd64 linux_amd64 linux_386 linux_arm)
files=(smartcooly-windows-4.0-amd64.exe smartcooly-windows-4.0-386.exe smartcooly-darwin-10.6-amd64 smartcooly-linux-amd64 smartcooly-linux-386 smartcooly-linux-arm-5)

unzip web/dist.zip -d web

for i in 0 1 2 3 4 5; do
  mkdir cache/smartcooly_${osarchs[${i}]}
  mkdir cache/smartcooly_${osarchs[${i}]}/web
  mkdir cache/smartcooly_${osarchs[${i}]}/custom
  cp LICENSE cache/smartcooly_${osarchs[${i}]}/LICENSE
  cp -r plugin cache/smartcooly_${osarchs[${i}]}/plugin
  cp README.md cache/smartcooly_${osarchs[${i}]}/README.md
  cp -r web/dist cache/smartcooly_${osarchs[${i}]}/web/dist
  cp config.ini cache/smartcooly_${osarchs[${i}]}/custom/config.ini
  cp config.ini cache/smartcooly_${osarchs[${i}]}/custom/config.default.ini
  cd cache
  if [ ${i} -lt 2 ]
  then
    mv ${files[${i}]} smartcooly_${osarchs[${i}]}/smartcooly.exe
    zip -r smartcooly_${osarchs[${i}]}.zip smartcooly_${osarchs[${i}]}
  else
    mv ${files[${i}]} smartcooly_${osarchs[${i}]}/smartcooly
    tar -zcvf smartcooly_${osarchs[${i}]}.tar.gz smartcooly_${osarchs[${i}]}
  fi
  rm -rf smartcooly_${osarchs[${i}]}
  cd ..
done

zip -r ./cache.zip ./cache/

rm -rf web/dist cache
