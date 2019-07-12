package collector

// import (
// 	"fmt"
// 	"net/http"
// 	_ "net/http/pprof"
// 	"sort"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/prometheus/common/log"
// 	"github.com/prometheus/common/version"
// 	"github.com/prometheus/node_exporter/collector"
// )

// // handler wraps an unfiltered http.Handler but uses a filtered handler,
// // created on the fly, if filtering is requested. Create instances with
// // newHandler.
// type handler struct {
// 	unfilteredHandler http.Handler
// 	// exporterMetricsRegistry is a separate registry for the metrics about
// 	// the exporter itself.
// 	exporterMetricsRegistry *prometheus.Registry
// 	includeExporterMetrics  bool
// 	maxRequests             int
// }

// func NewHandler(includeExporterMetrics bool, maxRequests int) *handler {
// 	h := &handler{
// 		exporterMetricsRegistry: prometheus.NewRegistry(),
// 		includeExporterMetrics:  includeExporterMetrics,
// 		maxRequests:             maxRequests,
// 	}
// 	if h.includeExporterMetrics {
// 		h.exporterMetricsRegistry.MustRegister(
// 			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
// 			prometheus.NewGoCollector(),
// 		)
// 	}
// 	if innerHandler, err := h.innerHandler(); err != nil {
// 		log.Fatalf("Couldn't create metrics handler: %s", err)
// 	} else {
// 		h.unfilteredHandler = innerHandler
// 	}
// 	return h
// }

// // ServeHTTP implements http.Handler.
// func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	filters := r.URL.Query()["collect[]"]
// 	log.Debugln("collect query:", filters)

// 	if len(filters) == 0 {
// 		// No filters, use the prepared unfiltered handler.
// 		h.unfilteredHandler.ServeHTTP(w, r)
// 		return
// 	}
// 	// To serve filtered metrics, we create a filtering handler on the fly.
// 	filteredHandler, err := h.innerHandler(filters...)
// 	if err != nil {
// 		log.Warnln("Couldn't create filtered metrics handler:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
// 		return
// 	}
// 	filteredHandler.ServeHTTP(w, r)
// }

// // innerHandler is used to create buth the one unfiltered http.Handler to be
// // wrapped by the outer handler and also the filtered handlers created on the
// // fly. The former is accomplished by calling innerHandler without any arguments
// // (in which case it will log all the collectors enabled via command-line
// // flags).
// func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
// 	nc, err := collector.NewNodeCollector(filters...)
// 	if err != nil {
// 		return nil, fmt.Errorf("couldn't create collector: %s", err)
// 	}

// 	// Only log the creation of an unfiltered handler, which should happen
// 	// only once upon startup.
// 	if len(filters) == 0 {
// 		log.Infof("Enabled collectors:")
// 		collectors := []string{}
// 		for n := range nc.Collectors {
// 			collectors = append(collectors, n)
// 		}
// 		sort.Strings(collectors)
// 		for _, n := range collectors {
// 			log.Infof(" - %s", n)
// 		}
// 	}

// 	r := prometheus.NewRegistry()
// 	r.MustRegister(version.NewCollector("node_exporter"))
// 	if err := r.Register(nc); err != nil {
// 		return nil, fmt.Errorf("couldn't register node collector: %s", err)
// 	}
// 	handler := promhttp.HandlerFor(
// 		prometheus.Gatherers{h.exporterMetricsRegistry, r},
// 		promhttp.HandlerOpts{
// 			ErrorLog:            log.NewErrorLogger(),
// 			ErrorHandling:       promhttp.ContinueOnError,
// 			MaxRequestsInFlight: h.maxRequests,
// 		},
// 	)
// 	if h.includeExporterMetrics {
// 		// Note that we have to use h.exporterMetricsRegistry here to
// 		// use the same promhttp metrics for all expositions.
// 		handler = promhttp.InstrumentMetricHandler(
// 			h.exporterMetricsRegistry, handler,
// 		)
// 	}
// 	return handler, nil
// }
