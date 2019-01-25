default:
	cd 6_Harness/Aux && go build -o '../xboard'
	cd ../../

	cd 7_Harness/Aux && go build -o '../xrules'
	cd ../../

	cd 8_Harness/Aux && go build -o '../xstrategy'
	cd ../../

	cd 10_Harness/Aux && go build -o '../xobserve'
	cd ../../

	cd 11_Harness/Aux && go build -o '../xrun'
	cd ../../

	cd 13_Client/Aux && go build -o '../xclients'
	cd ../../

	cd 13_Server/Aux && go build -o '../xserver'
	cd ../../

	# cd Player && make

copy:
	cp 6_Harness/xboard ../6/xboard
	cp 7_Harness/xrules ../7/xrules
	cp 8_Harness/xstrategy ../8/xstrategy
	mv 10_Harness/xobserve ../10/xobserve
	mv 11_Harness/xrun ../11/xrun
	mv 13_Client/xclients ../13/xclients
	mv 13_Server/xserver ../13/xserver

linux:
	cd 6_Harness/Aux && GOOS=linux go build -o '../xboard'
	cd ../../

	cd 7_Harness/Aux && GOOS=linux go build -o '../xrules'
	cd ../../

	cd 8_Harness/Aux && GOOS=linux go build -o '../xstrategy'
	cd ../../

	cd 10_Harness/Aux && GOOS=linux go build -o '../xobserve'
	cd ../../

	cd 11_Harness/Aux && GOOS=linux go build -o '../xrun'
	cd ../../

	cd 13_Client/Aux && GOOS=linux go build -o '../xclients'
	cd ../../

	cd 13_Server/Aux && GOOS=linux go build -o '../xserver'
	cd ../../

	# cd Player && make linux
