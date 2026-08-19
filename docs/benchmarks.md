# Benchmark Results

Performance benchmarks for the Go Dependency Injector.

> **Auto-generated** from CI on $(date -u +"%Y-%m-%d %H:%M UTC")
>
> Runner: GitHub Actions (ubuntu-latest)

## Summary

| Operation | Performance |
|-----------|-------------|
| Container creation | **33.65 ns** |
| Instance resolution | **69.96 ns** |
| Singleton resolution | **92.45 ns** |
| Registration | **105.7 ns** |

---

## Registration Operations

These benchmarks measure the cost of registering dependencies with the container.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkNew` | 33.65 | 0 | 0 |
| `BenchmarkNew` | 33.61 | 0 | 0 |
| `BenchmarkNew` | 34.55 | 0 | 0 |
| `BenchmarkRegister` | 105.7 | 96 | 1 |

---

## Resolution Operations

These benchmarks measure the cost of resolving dependencies from the container.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveTransient` | 316.5 | 56 | 3 |
| `BenchmarkResolveTransient` | 313.9 | 56 | 3 |
| `BenchmarkResolveTransient` | 314.9 | 56 | 3 |
| `BenchmarkResolveSingleton` | 92.45 | 16 | 1 |
| `BenchmarkResolveSingleton` | 93.51 | 16 | 1 |
| `BenchmarkResolveSingleton` | 92.72 | 16 | 1 |

---

## Dependency Chain Resolution

These benchmarks measure resolution performance with dependency injection chains.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveWithOneDependency` | 685.1 | 144 | 7 |
| `BenchmarkResolveWithOneDependency` | 671.0 | 144 | 7 |
| `BenchmarkResolveWithOneDependency` | 673.9 | 144 | 7 |

---

## Parallel/Concurrent Performance

These benchmarks measure thread-safe concurrent resolution.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveSingletonParallel` | 127.0 | 16 | 1 |
| `BenchmarkResolveSingletonParallel` | 120.3 | 16 | 1 |
| `BenchmarkResolveSingletonParallel` | 119.4 | 16 | 1 |
| `BenchmarkResolveTransientParallel` | 190.9 | 56 | 3 |

---

## Utility Operations

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkHas` | 29.04 | 0 | 0 |
| `BenchmarkHas` | 29.00 | 0 | 0 |
| `BenchmarkHas` | 29.00 | 0 | 0 |

---

## Large Container Performance

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkContainerWithManyRegistrations` | 10236 | 10208 | 106 |
| `BenchmarkContainerWithManyRegistrations` | 10455 | 10208 | 106 |

---

## Recommendations

Based on these benchmarks:

1. **Use Singletons** for stateless services — significantly faster than transients after initial creation
2. **Use RegisterInstance** for pre-created objects — fastest resolution path
3. **Use Scoped** for request-scoped dependencies — excellent parallel performance
4. **Minimize transient chains** — each level adds allocation overhead

---

## Running Benchmarks Locally

```bash
# Run all benchmarks
go test ./di/... -bench=. -benchmem

# Run specific benchmark
go test ./di/... -bench=BenchmarkResolveSingleton -benchmem

# Run with multiple iterations for accuracy
go test ./di/... -bench=. -benchmem -count=5
```

---

## Raw Output

<details>
<summary>Click to expand full benchmark output</summary>

```
goos: linux
goarch: amd64
pkg: github.com/quinnjr/go-dependency-injector/di
cpu: Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz
BenchmarkNew-4                              	35536258	        33.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkNew-4                              	35416683	        33.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkNew-4                              	32724375	        34.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkRegister-4                         	10897743	       105.7 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegister-4                         	11336530	       104.6 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegister-4                         	11544852	       101.0 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	11293436	       102.8 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	11795906	       100.2 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	11901540	       103.2 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	10921970	       110.4 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	10668844	       109.2 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	10372104	       109.6 ns/op	      96 B/op	       1 allocs/op
BenchmarkResolveTransient-4                 	 3745578	       316.5 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransient-4                 	 3744937	       313.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransient-4                 	 3793372	       314.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveSingleton-4                 	12857142	        92.45 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingleton-4                 	12805765	        93.51 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingleton-4                 	12740972	        92.72 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	16972737	        69.96 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	16949052	        70.31 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	16896799	        71.09 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	10113402	       114.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	10381976	       114.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	10400893	       114.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	12468412	        95.83 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	12492333	        94.98 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	11524438	        95.04 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveWithOneDependency-4         	 1793559	       685.1 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithOneDependency-4         	 1795224	       671.0 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithOneDependency-4         	 1776988	       673.9 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1247464	       960.7 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1279531	       953.5 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1261226	       967.2 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 1860550	       640.0 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 1875146	       641.5 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 1869398	       640.8 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveNamed-4                     	12721690	        93.80 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveNamed-4                     	12593740	        94.05 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveNamed-4                     	12471306	        94.40 ns/op	      16 B/op	       1 allocs/op
BenchmarkHas-4                              	41482202	        29.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkHas-4                              	41066661	        29.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkHas-4                              	41245188	        29.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	35963798	        33.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	35971896	        33.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	35766421	        33.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkCreateScope-4                      	11304183	       106.4 ns/op	     112 B/op	       2 allocs/op
BenchmarkCreateScope-4                      	11046651	       106.0 ns/op	     112 B/op	       2 allocs/op
BenchmarkCreateScope-4                      	11430621	       108.1 ns/op	     112 B/op	       2 allocs/op
BenchmarkResolveSingletonParallel-4         	 8769464	       127.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingletonParallel-4         	 9651663	       120.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingletonParallel-4         	10003473	       119.4 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveTransientParallel-4         	 6247790	       190.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransientParallel-4         	 6324387	       191.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransientParallel-4         	 6223171	       190.2 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveScopedParallel-4            	11803621	       103.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedParallel-4            	11657523	       103.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedParallel-4            	11266514	       104.9 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveWithDepsParallel-4          	 2993979	       400.3 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithDepsParallel-4          	 2985780	       398.1 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithDepsParallel-4          	 2991957	       394.5 ns/op	     144 B/op	       7 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  116230	     10236 ns/op	   10208 B/op	     106 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  109322	     10455 ns/op	   10208 B/op	     106 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  110062	     10922 ns/op	   10208 B/op	     106 allocs/op
BenchmarkResolveFromLargeContainer-4        	12628719	        94.38 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveFromLargeContainer-4        	12582852	        93.90 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveFromLargeContainer-4        	12884692	        92.59 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/quinnjr/go-dependency-injector/di	92.464s
```

</details>
