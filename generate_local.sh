#!/bin/bash

# Clean existing generated files
rm -rf gen/proto

# Create output directories
mkdir -p gen/proto/schema
mkdir -p gen/proto/service

# Generate schema files
protoc --proto_path=proto \
       --go_out=gen/proto \
       --go_opt=paths=source_relative \
       --go_opt=Mschema/user.proto=discord/gen/proto/schema \
       --go_opt=Mschema/channel.proto=discord/gen/proto/schema \
       --go_opt=Mschema/message.proto=discord/gen/proto/schema \
       --go_opt=Mschema/friend.proto=discord/gen/proto/schema \
       --go_opt=Mschema/text_channel.proto=discord/gen/proto/schema \
       --go_opt=Mschema/voice_channel.proto=discord/gen/proto/schema \
       --go_opt=Mschema/permission.proto=discord/gen/proto/schema \
       --go_opt=Mschema/sync.proto=discord/gen/proto/schema \
       proto/schema/*.proto

# Generate service files
for service in auth user channel message friend text_channel voice_channel permission; do
    protoc --proto_path=proto \
           --go_out=gen/proto \
           --go-grpc_out=gen/proto \
           --go_opt=paths=source_relative \
           --go-grpc_opt=paths=source_relative \
           --go_opt=Mschema/user.proto=discord/gen/proto/schema \
           --go_opt=Mschema/channel.proto=discord/gen/proto/schema \
           --go_opt=Mschema/message.proto=discord/gen/proto/schema \
           --go_opt=Mschema/friend.proto=discord/gen/proto/schema \
           --go_opt=Mschema/text_channel.proto=discord/gen/proto/schema \
           --go_opt=Mschema/voice_channel.proto=discord/gen/proto/schema \
           --go_opt=Mschema/permission.proto=discord/gen/proto/schema \
           --go_opt=Mschema/sync.proto=discord/gen/proto/schema \
           proto/service/$service/*.proto
done

echo "Proto generation completed!"
