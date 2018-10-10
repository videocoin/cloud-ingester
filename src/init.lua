require("string")
cjson = require("cjson")
prometheus = require("prometheus").init("prometheus_metrics")

metric_build_info = prometheus:gauge(
    "stream_ingest_build_info", 
    "Build info about stream ingest server",
    {"nginx_version", "nginx_rtmp_version", "compiler", "built"})

metric_uptime = prometheus:gauge(
    "stream_ingest_uptime",
    "Uptime of stream ingest server")

metric_clients = prometheus:gauge(
    "stream_ingest_clients_total",
    "Number of clients connected to the server",
    {"app", "stream", "state"})

metric_recording_streams = prometheus:gauge(
    "stream_ingest_recording_streams_total",
    "Number of recording streams",
    {"app", "stream"})

metric_active_streams = prometheus:gauge(
    "stream_ingest_active_streams_total",
    "Number of active streams",
    {"app", "stream"})

metric_publishing_streams = prometheus:gauge(
    "stream_ingest_publishing_streams_total",
    "Number of publishing streams",
    {"app", "stream"})

metric_showtime = prometheus:gauge(
    "stream_ingest_showtime_seconds",
    "Number of showtime seconds",
    {"app", "stream"})

metric_inbound_bytes = prometheus:gauge(
    "stream_ingest_inbound_bytes",
    "Number of inbound bytes",
    {"app", "stream"})

metric_outbound_bytes = prometheus:gauge(
    "stream_ingest_outbound_bytes",
    "Number of outbound bytes",
    {"app", "stream"})

metric_inbound_bandwidth = prometheus:gauge(
    "stream_ingest_inbound_bandwidth",
    "Number of inbound bandwidth",
    {"app", "stream"})

metric_outbound_bandwidth = prometheus:gauge(
    "stream_ingest_outbound_bandwidth",
    "Number of outbound bandwidth",
    {"app", "stream"})

function split(str, pat)
   local t = {}
   local fpat = "(.-)" .. pat
   local last_end = 1
   local s, e, cap = str:find(fpat, 1)
   while s do
      if s ~= 1 or cap ~= "" then
         table.insert(t,cap)
      end
      last_end = e+1
      s, e, cap = str:find(fpat, last_end)
   end
   if last_end <= #str then
      cap = str:sub(last_end)
      table.insert(t, cap)
   end
   return t
end

function split_path(str)
   return split(str,'[\\/]+')
end

function split_filename(str)
   return split(str,'-+')
end

function split_stream_id(str)
   return split(str,'-+')
end
