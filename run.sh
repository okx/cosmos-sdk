for pkg in $(go list ./... | grep -v staking | grep  -v distribution | grep -v gov | grep -v params); do
  go test -mod=readonly -tags='ledger test_ledger_mock' "$pkg"
done