wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3NTAwOTYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.eTSdPGHG8Oqx75wqwqG6-VskliosipOBqoCaFG5Iu9w"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3NTAwOTYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.eTSdPGHG8Oqx75wqwqG6-VskliosipOBqoCaFG5Iu9w; Path=/; Secure; HttpOnly; Expires=Sun, 26 May 2024 19:01:36 GMT;"
local counter = 1
request = function()
    local title = "title of my note " .. counter
    local content = "my text of my note " .. counter
    counter = counter + 1
    local body = string.format('{"data": {"title": "%s", "content": "%s"}}', title, content)
    return wrk.format(nil, "/api/note/add", nil, body)
end