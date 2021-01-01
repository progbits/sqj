#!/usr/bin/env bats

root=$(git rev-parse --show-toplevel)
sqj="${root}/bin/sqj"

@test "invalid invocation should exit with exit code 1" {
    run "${sqj}"
    [ "$status" -eq 1 ]
}

@test "select from basic array" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" 'SELECT id FROM []' -)

    [ $(echo "${result}" | sed -n '2p') == "5fef99445feb430d23e22be1" ]
    [ $(echo "${result}" | sed -n '4p') == "5fef9944c13168afef559442" ]
}

@test "select from basic array with condition" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" 'SELECT id FROM [] WHERE value < 28.5' -)

    [ $(echo "${result}" | sed -n '2p') == "5fef99449eb340d76e5abf84" ]
}
