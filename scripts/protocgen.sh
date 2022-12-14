#!/usr/bin/env bash

set -eo pipefail

# protoc_gen_gocosmos() {
#   if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
#     echo -e "\tPlease run this command from somewhere inside the ibc-go folder."
#     return 1
#   fi

#   go get github.com/cosmos/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
# }

# protoc_gen_gocosmos

cd proto
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # shellcheck disable=SC2044
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# move proto files to the right places
cp -r github.com/bianjieai/nft-transfer/* ./
rm -rf github.com
