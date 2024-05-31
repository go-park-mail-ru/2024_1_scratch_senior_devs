wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E"
wrk.headers["X-Csrf-Token"] = "2e495587-8347-4c0a-8041-d092df5a6ffa"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E; Path=/; Secure; HttpOnly; Expires=Mon, 29 May 2124 19:01:36 GMT;"

request = function()
    return wrk.format(nil, "/api/note/756d7250-1460-4314-9850-82da0f09cfff", nil, nil)
end
