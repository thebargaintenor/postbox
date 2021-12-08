POSTBOX_DIR := ~/.postbox

.PHONY: node
node:
	cd src/frontend; \
	nvm use; \
	npm i

.PHONY: client
client:
	cd src/frontend; \
	npm run bundle
	cp src/frontend/dist/bundle.js src/server/public

.PHONY: run
run:
	cd src/server; \
	go run postbox.go $(POSTBOX_DIR)
