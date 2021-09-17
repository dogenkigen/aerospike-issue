# Aerospike Issue

The project was created to present the Aerospike issue with using filtering statements while UDF execution. To see
available commands run`make help`. In order to run benchmarks presenting the issue run:

```bash
make run-benchmarks
```

The output will present time difference for filtered and not filtered UDF run. The benchmarks code can be found
in [repo_test.go](repository/repo_test.go) file.