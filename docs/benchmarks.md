# Benchmark Results

Performance benchmarks for the Go Dependency Injector.

> **Auto-generated** from CI on $(date -u +"%Y-%m-%d %H:%M UTC")
>
> Runner: GitHub Actions (ubuntu-latest)

## Summary

| Operation | Performance |
|-----------|-------------|
| Container creation | **28.35 ns** |
| Instance resolution | **52.28 ns** |
| Singleton resolution | **66.70 ns** |
| Registration | **79.24 ns** |

---

## Registration Operations

These benchmarks measure the cost of registering dependencies with the container.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkNew` | 28.35 | 0 | 0 |
| `BenchmarkNew` | 28.30 | 0 | 0 |
| `BenchmarkNew` | 28.44 | 0 | 0 |
| `BenchmarkRegister` | 79.24 | 96 | 1 |

---

## Resolution Operations

These benchmarks measure the cost of resolving dependencies from the container.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveTransient` | 205.9 | 56 | 3 |
| `BenchmarkResolveTransient` | 204.5 | 56 | 3 |
| `BenchmarkResolveTransient` | 206.2 | 56 | 3 |
| `BenchmarkResolveSingleton` | 66.70 | 16 | 1 |
| `BenchmarkResolveSingleton` | 67.06 | 16 | 1 |
| `BenchmarkResolveSingleton` | 66.46 | 16 | 1 |

---

## Dependency Chain Resolution

These benchmarks measure resolution performance with dependency injection chains.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveWithOneDependency` | 460.8 | 144 | 7 |
| `BenchmarkResolveWithOneDependency` | 463.8 | 144 | 7 |
| `BenchmarkResolveWithOneDependency` | 459.3 | 144 | 7 |

---

## Parallel/Concurrent Performance

These benchmarks measure thread-safe concurrent resolution.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkResolveSingletonParallel` | 54.94 | 16 | 1 |
| `BenchmarkResolveSingletonParallel` | 55.00 | 16 | 1 |
| `BenchmarkResolveSingletonParallel` | 54.58 | 16 | 1 |
| `BenchmarkResolveTransientParallel` | 103.8 | 56 | 3 |

---

## Utility Operations

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkHas` | 21.08 | 0 | 0 |
| `BenchmarkHas` | 21.09 | 0 | 0 |
| `BenchmarkHas` | 21.33 | 0 | 0 |

---

## Large Container Performance

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `BenchmarkContainerWithManyRegistrations` | 8375 | 10208 | 106 |
| `BenchmarkContainerWithManyRegistrations` | 8266 | 10208 | 106 |

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
cpu: AMD EPYC 9V74 80-Core Processor                
BenchmarkNew-4                              	41830144	        28.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkNew-4                              	42256042	        28.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkNew-4                              	36851748	        28.44 ns/op	       0 B/op	       0 allocs/op
BenchmarkRegister-4                         	15125766	        79.24 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegister-4                         	14788390	        78.84 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegister-4                         	14478757	        79.14 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	14851306	        80.46 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	15221427	        79.57 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterWithOptions-4              	14350288	        81.06 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	13898478	        87.58 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	13341662	        87.02 ns/op	      96 B/op	       1 allocs/op
BenchmarkRegisterInstance-4                 	14177545	        87.51 ns/op	      96 B/op	       1 allocs/op
BenchmarkResolveTransient-4                 	 5789185	       205.9 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransient-4                 	 5600611	       204.5 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransient-4                 	 5835776	       206.2 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveSingleton-4                 	17610535	        66.70 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingleton-4                 	17597026	        67.06 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingleton-4                 	17252401	        66.46 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	22455150	        52.28 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	21210208	        52.40 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveInstance-4                  	22532672	        52.20 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	13340830	        87.52 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	13426976	        87.40 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedSameScope-4           	13587022	        87.71 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	17342968	        68.18 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	17156043	        68.70 ns/op	      16 B/op	       1 allocs/op
BenchmarkMustResolve-4                      	17366538	        68.46 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveWithOneDependency-4         	 2603404	       460.8 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithOneDependency-4         	 2590622	       463.8 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithOneDependency-4         	 2612828	       459.3 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1822983	       662.3 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1812229	       658.4 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveWithTwoDependencies-4       	 1806978	       659.3 ns/op	     232 B/op	       9 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 2726060	       437.3 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 2741097	       441.4 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveDeepDependencyChain-4       	 2736675	       440.5 ns/op	     128 B/op	       6 allocs/op
BenchmarkResolveNamed-4                     	17275290	        67.87 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveNamed-4                     	17131572	        68.05 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveNamed-4                     	17336272	        67.96 ns/op	      16 B/op	       1 allocs/op
BenchmarkHas-4                              	56848592	        21.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkHas-4                              	56788329	        21.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkHas-4                              	56957017	        21.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	55106886	        23.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	55139041	        21.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasNamed-4                         	55194274	        22.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkCreateScope-4                      	14972019	        80.07 ns/op	     112 B/op	       2 allocs/op
BenchmarkCreateScope-4                      	15260568	        82.84 ns/op	     112 B/op	       2 allocs/op
BenchmarkCreateScope-4                      	15186703	        79.86 ns/op	     112 B/op	       2 allocs/op
BenchmarkResolveSingletonParallel-4         	21100425	        54.94 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingletonParallel-4         	21759267	        55.00 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveSingletonParallel-4         	20659323	        54.58 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveTransientParallel-4         	11252808	       103.8 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransientParallel-4         	10987522	       103.0 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveTransientParallel-4         	11482304	       103.1 ns/op	      56 B/op	       3 allocs/op
BenchmarkResolveScopedParallel-4            	24682066	        48.68 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedParallel-4            	20625865	        49.03 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveScopedParallel-4            	23870949	        57.22 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveWithDepsParallel-4          	 5284861	       228.5 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithDepsParallel-4          	 5266371	       229.3 ns/op	     144 B/op	       7 allocs/op
BenchmarkResolveWithDepsParallel-4          	 5199682	       229.5 ns/op	     144 B/op	       7 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  146400	      8375 ns/op	   10208 B/op	     106 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  140544	      8266 ns/op	   10208 B/op	     106 allocs/op
BenchmarkContainerWithManyRegistrations-4   	  146078	      8158 ns/op	   10208 B/op	     106 allocs/op
BenchmarkResolveFromLargeContainer-4        	17638954	        66.54 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveFromLargeContainer-4        	17405457	        67.27 ns/op	      16 B/op	       1 allocs/op
BenchmarkResolveFromLargeContainer-4        	17623183	        66.64 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/quinnjr/go-dependency-injector/di	88.943s
```

</details>
