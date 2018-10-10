local req_data = ngx.req.get_uri_args()

if req_data and req_data.call == "update_publish" then
    ngx.log(ngx.INFO, "hook=", req_data.call, " ", cjson.encode(req_data))

    local stream_parts = split_stream_id(req_data.name)
    local user_id = stream_parts[1]
    local camera_id = stream_parts[2]

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

    local res = ngx.location.capture(
        "/api/v1/cameras/" .. camera_id .. "/onair",
        { method = ngx.HTTP_GET }
    )

    if res then
        ngx.log(ngx.INFO, "cameras api returned status ", res.status)
        ngx.log(ngx.INFO, "cameras api returned body ", res.body)
        ngx.exit(res.status)
    end

end

ngx.exit(400)