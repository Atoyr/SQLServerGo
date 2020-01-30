all: 
	@make webbuild
	@make pack

webbuild:
	@cd ./webapps && npm run build
	@rm -rf public
	@mkdir public
	@cp -r ./webapps/dist/. ./public

pack:
	packr2 clean
	packr2
