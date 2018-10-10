local res = ngx.location.capture("/stat.json")
if res.status >= 500 then 
    ngx.exit(res.status) 
end

ngx.status = res.status

local stat = cjson.decode(res.body)

metric_build_info:set(1, {
  stat.rtmp.nginx_version, 
  stat.rtmp.nginx_rtmp_version, 
  stat.rtmp.compiler:match("^%s*(.-)%s*$"), 
  stat.rtmp.built
})

metric_uptime:set(stat.rtmp.uptime)

local publishing = "publishing"
local playing = "playing"

if stat.rtmp.servers then
  for _, servers in pairs(stat.rtmp.servers) do
    for _, app in pairs(servers) do
        if app.live then
          for _, stream in pairs(app.live.streams) do
            local publishing_count = 0
            local playing_count = 0
            for _, client in pairs(stream.clients) do
              if client.publishing then
                publishing_count = publishing_count + 1
              else
                playing_count = playing_count + 1
              end
            end
            metric_clients:set(publishing_count, {app.name, stream.name, publishing})
            metric_clients:set(playing_count, {app.name, stream.name, playing})

            local recording_count = 0
            if stream.recording then
              recording = 1
            end
            metric_recording_streams:set(recording_count, {app.name, stream.name})

            local publishing_count = 0
            if stream.publishing then
              publishing_count = 1
            end
            metric_publishing_streams:set(publishing_count, {app.name, stream.name})

            local active_count = 0
            if stream.active then
              active_count = 1
            end
            metric_active_streams:set(active_count, {app.name, stream.name})

            metric_showtime:set(stream.time / 1000, {app.name, stream.name})

            metric_inbound_bytes:set(stream.bytes_in, {app.name, stream.name})
            metric_outbound_bytes:set(stream.bytes_out, {app.name, stream.name})

            metric_inbound_bandwidth:set(stream.bw_in, {app.name, stream.name})
            metric_outbound_bandwidth:set(stream.bw_out, {app.name, stream.name})
          end
        end
    end
  end
end

prometheus:collect()
