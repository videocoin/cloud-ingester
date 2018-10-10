local req_data = ngx.req.get_uri_args()

if req_data and req_data.call == "record_done" then
    ngx.log(ngx.INFO, "system=system hook=", req_data.call, " ", cjson.encode(req_data))

    local camera_id = req_data.name
    local path = req_data.path
    local recorder = req_data.recorder
    local parts = split_path(path)
    local filename_parts = split_filename(parts[3])
    local user_id = filename_parts[1]
    local record_file = parts[3]

    if recorder == "stream" then
        local res_user = ngx.location.capture(
            "/api/private/v1/users/" .. user_id .. "/",
            { method = ngx.HTTP_GET }
        )

        if res_user.status ~= 200 then
            ngx.log(ngx.ERR, "system=users message=user not found status=", res_user.status)
            ngx.exit(400)
        end

        local user_info = cjson.decode(res_user.body)
        ngx.log(ngx.INFO, "system=users token=", user_info.token)

        ngx.req.set_header("Authorization", "Bearer " .. user_info.token)
        ngx.req.set_header("Content-Type", "application/json")

        local url = "http://" .. ngx.var.pod_ip .. "/records/" .. parts[3]

        ngx.log(ngx.INFO, "url=", url)

        local res = ngx.location.capture(
            "/api/v1/records",
            {
                method = ngx.HTTP_POST,
                body = cjson.encode({
                    camera_id = camera_id,
                    url = url
                })
            }
        )

        if res then
            if res.status == ngx.HTTP_CREATED then
                ngx.log(ngx.INFO, "system=uploader message=record has been created")
            else
                ngx.log(ngx.INFO, "system=uploader message=failed to create record status=", res.status)
            end
            ngx.exit(res.status)
        end
    end

end

ngx.exit(400)