@echo off
set GO111MODULE=off
set GOPATH=C:\Users\WildClown\Documents\TCC\adaptive-20220814T014930Z-001\adaptive

rem compile server
cd C:\Users\WildClown\Documents\TCC\adaptive-20220814T014930Z-001\adaptive\src\rabbitmq\modified\img\server
go build -o server.exe main.go

rem go run main.go -is-adaptive=false -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=5000 -kp=0.005643669 -ki=-0.020420019 -kd=0.001412706
go run main.go -is-adaptive=false -monitor-interval=10 -prefetch-count=1 -controller-type="PID" -set-point=2000.0 -kp=-0.000325258 -ki=-0.011864766 -kd=0.000820833
echo FINISH
pause