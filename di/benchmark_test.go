package di_test

import (
	"testing"

	"github.com/quinnjr/go-dependency-injector/di"
)

// =============================================================================
// Benchmark Types
// =============================================================================

type BenchLogger interface {
	Log(msg string)
}

type benchLoggerImpl struct{}

func (l *benchLoggerImpl) Log(msg string) {}

type BenchService interface {
	DoWork() string
}

type benchServiceImpl struct {
	logger BenchLogger
}

func (s *benchServiceImpl) DoWork() string {
	return "done"
}

type BenchComplexService interface {
	Process() string
}

type benchComplexServiceImpl struct {
	logger  BenchLogger
	service BenchService
}

func (s *benchComplexServiceImpl) Process() string {
	s.logger.Log("processing")
	return s.service.DoWork()
}

// =============================================================================
// Registration Benchmarks
// =============================================================================

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = di.New()
	}
}

func BenchmarkRegister(b *testing.B) {
	c := di.New()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		di.Register[BenchLogger](c, func() BenchLogger {
			return &benchLoggerImpl{}
		})
	}
}

func BenchmarkRegisterWithOptions(b *testing.B) {
	c := di.New()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		di.Register[BenchLogger](c, func() BenchLogger {
			return &benchLoggerImpl{}
		}, di.AsSingleton())
	}
}

func BenchmarkRegisterInstance(b *testing.B) {
	c := di.New()
	instance := &benchLoggerImpl{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		di.RegisterInstance[BenchLogger](c, instance)
	}
}

// =============================================================================
// Resolution Benchmarks
// =============================================================================

func BenchmarkResolveTransient(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchLogger](c)
	}
}

func BenchmarkResolveSingleton(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	// Warm up singleton
	_, _ = di.Resolve[BenchLogger](c)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchLogger](c)
	}
}

func BenchmarkResolveInstance(b *testing.B) {
	c := di.New()
	di.RegisterInstance[BenchLogger](c, &benchLoggerImpl{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchLogger](c)
	}
}

func BenchmarkResolveScopedSameScope(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsScoped())

	scope := c.CreateScope("bench")

	// Warm up scope
	_, _ = di.ResolveInScope[BenchLogger](c, scope)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.ResolveInScope[BenchLogger](c, scope)
	}
}

func BenchmarkMustResolve(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	// Warm up
	_ = di.MustResolve[BenchLogger](c)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = di.MustResolve[BenchLogger](c)
	}
}

// =============================================================================
// Dependency Chain Benchmarks
// =============================================================================

func BenchmarkResolveWithOneDependency(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	di.Register[BenchService](c, func(l BenchLogger) BenchService {
		return &benchServiceImpl{logger: l}
	}, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchService](c)
	}
}

func BenchmarkResolveWithTwoDependencies(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	di.Register[BenchService](c, func(l BenchLogger) BenchService {
		return &benchServiceImpl{logger: l}
	}, di.AsSingleton())

	di.Register[BenchComplexService](c, func(l BenchLogger, s BenchService) BenchComplexService {
		return &benchComplexServiceImpl{logger: l, service: s}
	}, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchComplexService](c)
	}
}

func BenchmarkResolveDeepDependencyChain(b *testing.B) {
	c := di.New()

	// Create a chain: Level5 -> Level4 -> Level3 -> Level2 -> Level1
	type Level1 interface{ L1() }
	type Level2 interface{ L2() }
	type Level3 interface{ L3() }
	type Level4 interface{ L4() }
	type Level5 interface{ L5() }

	di.Register[Level1](c, func() Level1 { return &level1Impl{} }, di.AsSingleton())
	di.Register[Level2](c, func(l1 Level1) Level2 { return &level2Impl{} }, di.AsSingleton())
	di.Register[Level3](c, func(l2 Level2) Level3 { return &level3Impl{} }, di.AsSingleton())
	di.Register[Level4](c, func(l3 Level3) Level4 { return &level4Impl{} }, di.AsSingleton())
	di.Register[Level5](c, func(l4 Level4) Level5 { return &level5Impl{} }, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[Level5](c)
	}
}

type level1Impl struct{}
type level2Impl struct{}
type level3Impl struct{}
type level4Impl struct{}
type level5Impl struct{}

func (l *level1Impl) L1() {}
func (l *level2Impl) L2() {}
func (l *level3Impl) L3() {}
func (l *level4Impl) L4() {}
func (l *level5Impl) L5() {}

// =============================================================================
// Named Resolution Benchmarks
// =============================================================================

func BenchmarkResolveNamed(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.WithName("primary"), di.AsSingleton())

	// Warm up
	_, _ = di.ResolveNamed[BenchLogger](c, "primary")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.ResolveNamed[BenchLogger](c, "primary")
	}
}

// =============================================================================
// Utility Benchmarks
// =============================================================================

func BenchmarkHas(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = di.Has[BenchLogger](c)
	}
}

func BenchmarkHasNamed(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.WithName("named"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = di.HasNamed[BenchLogger](c, "named")
	}
}

func BenchmarkCreateScope(b *testing.B) {
	c := di.New()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.CreateScope("scope")
	}
}

// =============================================================================
// Concurrent Benchmarks
// =============================================================================

func BenchmarkResolveSingletonParallel(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	// Warm up
	_, _ = di.Resolve[BenchLogger](c)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.Resolve[BenchLogger](c)
		}
	})
}

func BenchmarkResolveTransientParallel(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.Resolve[BenchLogger](c)
		}
	})
}

func BenchmarkResolveScopedParallel(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsScoped())

	scope := c.CreateScope("parallel")

	// Warm up
	_, _ = di.ResolveInScope[BenchLogger](c, scope)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.ResolveInScope[BenchLogger](c, scope)
		}
	})
}

func BenchmarkResolveWithDepsParallel(b *testing.B) {
	c := di.New()
	di.Register[BenchLogger](c, func() BenchLogger {
		return &benchLoggerImpl{}
	}, di.AsSingleton())

	di.Register[BenchService](c, func(l BenchLogger) BenchService {
		return &benchServiceImpl{logger: l}
	}, di.AsTransient())

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.Resolve[BenchService](c)
		}
	})
}

// =============================================================================
// Memory Benchmarks
// =============================================================================

func BenchmarkContainerWithManyRegistrations(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		c := di.New()

		// Register 100 types
		for j := 0; j < 100; j++ {
			di.Register[BenchLogger](c, func() BenchLogger {
				return &benchLoggerImpl{}
			})
		}
	}
}

func BenchmarkResolveFromLargeContainer(b *testing.B) {
	c := di.New()

	// Register many types to simulate real-world container
	for i := 0; i < 100; i++ {
		di.Register[BenchLogger](c, func() BenchLogger {
			return &benchLoggerImpl{}
		}, di.AsSingleton())
	}

	// Warm up
	_, _ = di.Resolve[BenchLogger](c)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[BenchLogger](c)
	}
}
