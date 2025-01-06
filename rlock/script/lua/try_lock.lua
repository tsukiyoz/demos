local val = redis.call('GET', KEYS[1])

if val ~= false then
    return 0
end


local res = redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2]) 
if res['ok'] == "OK" then
    return 1
else
    return 0
end
