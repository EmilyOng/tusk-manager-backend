# Builds the application
build:
	go build -o main .

# Start up a development server with live-reload utility
start:
	$(shell go env GOPATH)/bin/air

# Generate types to be used on frontend
generate-types:
	rm -rf ../tusk-manager-frontend/src/generated
	mkdir ../tusk-manager-frontend/src/generated
	touch ../tusk-manager-frontend/src/generated/types.ts
	# Handle Enums in types/color and types/role
	echo "export enum Color {Turquoise = 'Turquoise', Blue = 'Blue', Cyan = 'Cyan', Green = 'Green', Yellow = 'Yellow', Red = 'Red'}" >> ../tusk-manager-frontend/src/generated/types.ts 
	echo "export enum Role {Owner = 'Owner', Editor = 'Editor', Viewer = 'Viewer'}" >> ../tusk-manager-frontend/src/generated/types.ts 
	touch ../tusk-manager-frontend/src/generated/views.ts
	$(shell go env GOPATH)/bin/tscriptify \
		-package=github.com/EmilyOng/tusk-manager/backend/views \
		-target=../tusk-manager-frontend/src/generated/views.ts \
		-import="import { Color } from './types'" \
		-import="import { Role } from './types'" \
		-interface \
		views/auth.go \
		views/board.go \
		views/member.go \
		views/response.go \
		views/state.go \
		views/tag.go \
		views/task.go \
		views/user.go
