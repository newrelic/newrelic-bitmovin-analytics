# Build, deploy lambda, and clean all. Must specify lambda names in command line.
#
# Usage:
#    make recv=ReceiverLambda proc=ProcessorLambda expt=ExporterLambda

arch 	= amd64

# Lambda names

recv	= ""
proc	= ""
expt	= ""

deploy: deploy-recv deploy-proc deploy-expt

# Receiver part

deploy-recv: build-recv update-recv clean-recv

compile-recv:
	( cd ./receiver && GOOS=linux GOARCH=$(arch) go build -o receiver )

zip-recv:
	( cd ./receiver && zip receiver.zip receiver )

update-recv:
	( cd ./receiver && aws lambda update-function-code --no-cli-pager --function-name $(recv) --zip-file fileb://receiver.zip )

clean-bin-recv:
	( cd ./receiver && rm -f receiver )

clean-zip-recv:
	( cd ./receiver && rm -f receiver.zip )

clean-recv: clean-bin-recv clean-zip-recv

build-recv: compile-recv zip-recv

# Processor part

deploy-proc: build-proc update-proc clean-proc

compile-proc:
	( cd ./processor && GOOS=linux GOARCH=$(arch) go build -o processor )

zip-proc:
	( cd ./processor && zip processor.zip processor )

update-proc:
	( cd ./processor && aws lambda update-function-code --no-cli-pager --function-name $(proc) --zip-file fileb://processor.zip )

clean-bin-proc:
	( cd ./processor && rm -f processor )

clean-zip-proc:
	( cd ./processor && rm -f processor.zip )

clean-proc: clean-bin-proc clean-zip-proc

build-proc: compile-proc zip-proc

# Exporter part

deploy-expt: build-expt update-expt clean-expt

compile-expt:
	( cd ./exporter && GOOS=linux GOARCH=$(arch) go build -o exporter )

zip-expt:
	( cd ./exporter && zip exporter.zip exporter )

update-expt:
	( cd ./exporter && aws lambda update-function-code --no-cli-pager --function-name $(expt) --zip-file fileb://exporter.zip )

clean-bin-expt:
	( cd ./exporter && rm -f exporter )

clean-zip-expt:
	( cd ./exporter && rm -f exporter.zip )

clean-expt: clean-bin-expt clean-zip-expt

build-expt: compile-expt zip-expt
