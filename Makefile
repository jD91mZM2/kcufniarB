default:
	go install
build:
	# Won't work for anybody but me.
	# Requires a special script.
	cd $(GOPATH)/bin; \
		./Cross\ Compile.sh kcufniarB
