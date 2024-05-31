wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E"
wrk.headers["X-Csrf-Token"] = "2e495587-8347-4c0a-8041-d092df5a6ffa"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E; Path=/; Secure; HttpOnly; Expires=Mon, 29 May 2124 19:01:36 GMT;"

local counter = 0

request = function()
    counter = counter + 1

    local title = "title of my note " .. counter
    local content = "my text of my note " .. counter
    local body = string.format('{"data": {"title": "%s", "content": "%s"}}', title, content)
    return wrk.format(nil, "/api/note/add", nil, body)
end