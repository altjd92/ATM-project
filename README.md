# ATM-project
1. SmartContract 생성
  go mod init 생성
  go mod vendor 폴더 생성
  go build

2. 네트워크 구조
  1) 채널 생성
  ./network.sh up createChannel -ca
  2) peer 생성
  ./network.sh deployCC -ccn atm -ccp ../contract/ -ccl go -ccv 1.0
  
 3. application 실행
  1) npm 설치
  npm install
  2) main.js 실행
  
