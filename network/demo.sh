#!/bin/bash

# 네트워크 시작
peer chaincodeinvoke -o localhost:7050 --ordererTLSHostnameOverrideorderer.example.com --tls--cafile$ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"InitLedger","Args":[]}'
# 체인코드 배포
peer lifecycle chaincodecommit -o localhost:7050 --ordererTLSHostnameOverrideorderer.example.com --tls--cafile$ORDERER_CA --channelID$CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence 1


