
proto:
	cd common/proto;protoc --proto_path=. --proto_path=/Users/shicongying/Downloads/Goland_Work/src --micro_out=. --gogofaster_out=plugins=grpc:. *.proto
	cd common/proto;grep -rl ',omitempty' ./*.go | xargs sed -i '' 's/,omitempty//g'