#!/usr/bin/env bats

root=$(git rev-parse --show-toplevel)
sqj="${root}/bin/sqj"

@test "invalid invocation should exit with EXIT_FAILURE" {
    run "${sqj}"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Usage"* ]]
}

@test "select from an empty array" {
    result=$(echo "[]" | "${sqj}" --compact 'SELECT * FROM []' -)
    expected='[]'
    [ "${result}" == "${expected}" ]
}

@test "select all columns from basic array" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" --compact 'SELECT * FROM []' -)
    expected='[{"id":"5fef99445feb430d23e22be1","value":47.5088},{"id":"5fef99449eb340d76e5abf84","value":28.4475},{"id":"5fef9944c13168afef559442","value":49.9649}]'
    [ "${result}" == "${expected}" ]
}

@test "select specific column from basic array" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" --compact 'SELECT value FROM []' -)
    expected='[{"value":47.5088},{"value":28.4475},{"value":49.9649}]'
    [ "${result}" == "${expected}" ]
}

@test "select specific column basic array with condition" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" --compact 'SELECT id FROM [] WHERE value < 28.5' -)
    expected='[{"id":"5fef99449eb340d76e5abf84"}]'
    [ "${result}" == "${expected}" ]
}

@test "select aliased member of a nested object with a nested object based condition" {
    result=$(cat "${root}/test/data/array_with_nested.json" | "${sqj}" --compact 'SELECT about$registered AS date FROM [] WHERE about$metric > 0' -)
    expected='[{"date":"2016-03-21T09:47:51 -00:00"},{"date":"2015-12-17T09:14:19 -00:00"}]'
    [ "${result}" == "${expected}" ]
}

@test "select aliased member of a nested object with a nested object based condition" {
    result=$(cat "${root}/test/data/array_with_nested.json" | "${sqj}" --compact 'SELECT about$registered AS date FROM [] WHERE about$metric > 0' -)
    expected='[{"date":"2016-03-21T09:47:51 -00:00"},{"date":"2015-12-17T09:14:19 -00:00"}]'
    [ "${result}" == "${expected}" ]
}

