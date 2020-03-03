#!/usr/bin/env bash

for i in {0..100..1}
do
   echo "play n°$i"
   go test -v -race -failfast --cover -covermode=atomic -mod=vendor
   if [ $? -ne 0 ]; then
      echo -e "failed at play n° $i"
      exit $?
   fi
   sleep 5
done
