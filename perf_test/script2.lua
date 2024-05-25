wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3NDUzMTcsImlkIjoiNmIxMGFmNzgtZTk0Yy00YmJhLTg2NDMtMjc4ZjIxZmY1MTYzIiwidXNyIjoidGVzdHVzZXIifQ.2BKPCQ6foH6rc_TL59FlAdK7htBk9j3l3XB1jMfYYCw"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3NDUzMTcsImlkIjoiNmIxMGFmNzgtZTk0Yy00YmJhLTg2NDMtMjc4ZjIxZmY1MTYzIiwidXNyIjoidGVzdHVzZXIifQ.2BKPCQ6foH6rc_TL59FlAdK7htBk9j3l3XB1jMfYYCw; Path=/; Secure; HttpOnly; Expires=Sun, 26 May 2024 17:41:57 GMT;YouNoteCSRF=4937b223-076b-4e7e-9da6-0e36bca185a3; Path=/; Expires=Sun, 26 May 2024 17:41:57 GMT; HttpOnly; Secure; SameSite=Strict"

local offset = 0

request = function()
    local limit = 10
    offset = offset + 10

    local body = string.format('/api/note/get_all?count=%s&offset=%s', limit, offset)

    return wrk.format(nil, body, nil)
end
