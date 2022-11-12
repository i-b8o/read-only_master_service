#!/bin/bash
echo "supreme update contracts"
cd app/
go get -u github.com/i-b8o/regulations_contracts@$1
