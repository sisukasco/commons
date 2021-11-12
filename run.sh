#!/bin/bash


function runAllTests(){
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user go test -p 1 -count=1 $(go list ./...)
}

function runTests(){
    echo "Running tests in folder $1 ..."
    go test -v -count=1 $1
}

function runTestFn(){
    echo "Running test $1 in folder $2 ..."
    go test -v -count=1 -run $1 $2
}


if (( $# < 2 ))
then
      echo "Usage cmd.sh command what"
      exit 1
fi

COMMAND=$1
WHAT=$2

case $COMMAND in
    test)
        case $WHAT in
           local)
                runAllTests
           ;;
           fn)
                runTestFn $3 $4
           ;; 
           *)
                runTests $WHAT
           ;;
        esac
        
    ;;
esac

