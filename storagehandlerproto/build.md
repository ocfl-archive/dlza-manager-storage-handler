
### Generate GRPC code for dlza manager storage handler
```bash
protoc --proto_path=. --proto_path=C:\Users\jarek\GolandProjects\github\dlza-manager\dlzamanagerproto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
      