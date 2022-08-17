@echo off
set GO111MODULE=off
set currentpath=%~dp0
echo %currentpath:~0,-1%
cd..\..\..\..\..
echo %cd%
set GOPATH=%cd%
rem compile server
cd %currentpath:~0,-1%
go build -o server.exe main.go

rem go run main.go -is-adaptive=false -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=5000 -kp=0.005643669 -ki=-0.020420019 -kd=0.001412706
go run main.go -is-adaptive=false -monitor-interval=3 -prefetch-count=1 -controller-type="PID" -set-point=2000.0 -kp=-0.000325258 -ki=-0.011864766 -kd=0.000820833
echo FINISH
pause