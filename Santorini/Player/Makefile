default:
	cd Valid && go build -o "../../plugins/valid.so" -buildmode=plugin
	cd ../

	cd Broken && go build -o "../../plugins/broken.so" -buildmode=plugin
	cd ../

	cd InfPlace && go build -o "../../plugins/infplace.so" -buildmode=plugin
	cd ../

	cd InfTurn && go build -o "../../plugins/infturn.so" -buildmode=plugin
	cd ../

linux:
	cd Valid && GOOS=linux go build -o "../../plugins/valid.so" -buildmode=plugin
	cd ../

	cd Broken && GOOS=linux go build -o "../../plugins/broken.so" -buildmode=plugin
	cd ../

	cd InfPlace && GOOS=linux go build -o "../../plugins/infplace.so" -buildmode=plugin
	cd ../

	cd InfTurn && GOOS=linux go build -o "../../plugins/infturn.so" -buildmode=plugin
	cd ../