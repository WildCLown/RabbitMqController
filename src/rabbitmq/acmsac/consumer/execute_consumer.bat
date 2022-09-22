@echo off
set GO111MODULE=off
set currentpath=%~dp0
echo %currentpath:~0,-1%
cd..\..\..\..
echo %cd%
set GOPATH=%cd%

cd %currentpath:~0,-1%

echo 1)Remember to start rabbimq-server with 'brew services start rabbitmq'
echo OR Stop 'brew services stop rabbitmq
echo 2)Remeber that PC=0 is infinite buffer

echo Compiling main.go
go build main.go

echo Server started...

@REM  Training || Change interval to 10
@REM go run main.go -is-adaptive=false -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=0 -kp=0.0 -ki=0.0 -kd=0.0 -csv-printer=true

@REM 'Root Locus'
@REM PI
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1400.0 -kp=-0.0003 -ki=0.0025 -kd=0.0 -csv-printer=true -label="RL-PI"

@REM PID
go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1400.0 -kp=0.0010903 -ki=0.0024250 -kd=0.0007687 -csv-printer=true -label="RL-PID"

@REM Ziegler-Nichols
@REM PI
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1400.0 -kp=0.001192 -ki=0.00003 -kd=0.0 -csv-printer=true -label="ZN-PI"
@REM PID
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1400.0 -kp=0.001589 -ki=0.00002 -kd=0.000005 -csv-printer=true -label="ZN-PID"



@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1000 -kp=0.0179 -ki=0.0 -kd=0.0
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=1000 -kp=0.166577703 -ki=0.158808685 -kd=0.036932252

@REM P Controller [1-22]
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.021147 -ki=0.0 -kd=0.0

@REM PI Controller [1-22] VEM
@REM echo [ PI Controller - Analytical ]
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=-0.0016 -ki=0.0060 -kd=0.0

@REM  PI - Error Square
@REM echo [ PI Error Square - Analytical ]
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="ErrorSquare" -set-point=400 -kp=-0.0016 -ki=0.0060 -kd=0.0

@REM PI controller [1-22] ACM SAC
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=-0.0016 -ki=0.0060 -kd=0.0

@REM PI Controller - Ziegler-Nichols
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.015475 -ki=51.583617 -kd=0.0

@REM  PI Controller - Cohen
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.059876 -ki=532.2303598 -kd=0.0

@REM  PI Controller - AMIGO
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.0086626 -ki=86.62603513 -kd=0.0

@REM  PD Controller [1-22] ACM SAC
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.09 -ki=0.0 -kd=0.09

@REM  PID [1-22] VEM
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.0017761 -ki=0.0058096 -kd=0.0018417

@REM  PID [1-22] ACM SAC
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.0017761 -ki=0.0058096 -kd=0.0018417

@REM  PID [1-22] Cohen
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.014968 -ki=80.32624361 -kd=0.00000046538

@REM  PID [1-22] AMIGO
@REM go run main.go -is-adaptive=true -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=400 -kp=0.022523 -ki=206.4587171 -kd=0.00000086626

@REM  OnOff
@REM go run main.go -is-adaptive=true -monitor-interval=30 -prefetch-count=1 -controller-type="OnOff" -set-point=400 -kp=0.0 -ki=0.0 -kd=0.0

@REM  OnOffDeadZone
@REM go run main.go -is-adaptive=true -monitor-interval=5 -prefetch-count=1 -controller-type="OnOffDeadZone" -set-point=400 -kp=0.0 -ki=0.0 -kd=0.0
echo FINISH
pause