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

@test "select from an empty object" {
    result=$(echo "{}" | "${sqj}" --compact 'SELECT * FROM []' -)
    expected='{}'
    [ "${result}" == "${expected}" ]
}

@test "input object with keyword members" {
  for input in "${root}"/test/data/keywords_*; do
    result=$(cat "${input}" | "${sqj}" --nth 0 'SELECT * FROM []' -)
    expected=$(cat "${input}")
    [ "${result}" == "${expected}" ]
  done
}

@test "select a null value" {
    result=$(echo '[{"name": null}]' | "${sqj}" --compact 'SELECT * FROM []' -)
    expected='[{"name":null}]'
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

@test "select all columns in the presence of nested objects" {
    result=$(cat "${root}/test/data/array_with_nested.json" | "${sqj}" --compact 'SELECT * FROM [] LIMIT 1' -)
    expected='[{"id":"5ff8d1fbdc7d0c09c7138193","guid":"dc640ea9-28f8-4ff6-8609-cfc6157e52bc","isActive":true,"about":{"registered":"2014-11-04T10:37:35 -00:00","metric":-24.0467}}]'
    [ "${result}" == "${expected}" ]
}

@test "select aliased member of a nested object with a nested object based condition" {
    result=$(cat "${root}/test/data/array_with_nested.json" | "${sqj}" --compact 'SELECT about$registered AS date FROM [] WHERE about$metric > 0' -)
    expected='[{"date":"2016-03-21T09:47:51 -00:00"},{"date":"2015-12-17T09:14:19 -00:00"}]'
    [ "${result}" == "${expected}" ]
}

@test "select based on JSON boolean condition" {
    result=$(cat "${root}/test/data/array_with_nested.json" | "${sqj}" --compact 'SELECT id FROM [] WHERE NOT isActive' -)
    expected='[{"id":"5ff8d1fbe962cc214df87658"},{"id":"5ff8d1fbd17bdee1c0768755"}]'
    [ "${result}" == "${expected}" ]
}

@test "select all columns from a JSON object" {
    result=$(cat "${root}/test/data/object.json" | "${sqj}" --compact 'SELECT id, guid, isActive FROM []' -)
    expected='[{"id":"5ff9d7a9380c2b3ffd350209","guid":"af6e799c-42d8-42e0-b11d-1ebbdc2f9f74","isActive":false}]'
    [ "${result}" == "${expected}" ]
}

@test "select specific columns from a JSON object" {
    result=$(cat "${root}/test/data/object.json" | "${sqj}" --compact 'SELECT isActive, guid FROM []' -)
    expected='[{"isActive":false,"guid":"af6e799c-42d8-42e0-b11d-1ebbdc2f9f74"}]'
    [ "${result}" == "${expected}" ]
}

@test "select nested array from a JSON object" {
    result=$(cat "${root}/test/data/object_with_nested_array.json" | "${sqj}" --compact 'SELECT friends FROM []' -)
    expected='[{"friends":[{"id":0,"name":"Dawson Joyce"},{"id":1,"name":"Trujillo Jones"},{"id":2,"name":"Sykes Cross"}]}]'
    [ "${result}" == "${expected}" ]
}

@test "--nth command line option - select first element" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" --compact --nth 0 'SELECT * FROM []' -)
    expected='{"id":"5fef99445feb430d23e22be1","value":47.5088}'
    [ "${result}" == "${expected}" ]
}

@test "--nth command line option" {
    result=$(cat "${root}/test/data/basic_array.json" | "${sqj}" --compact --nth 2 'SELECT * FROM []' -)
    expected='{"id":"5fef9944c13168afef559442","value":49.9649}'
    [ "${result}" == "${expected}" ]
}