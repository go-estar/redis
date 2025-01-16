package script

import "github.com/redis/go-redis/v9"

var BalanceGet = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* ARGV[1] 余额编号
*/]]
local val = redis.call('hget', KEYS[1], ARGV[1])
local total = 0
if val then 
	total = val
end
return total
`)

var BalanceList = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* ARGV[1] 余额编号
* ARGV[2] List[n]
*/]]
local ARGV2 = cjson.decode(ARGV[2])
local result = {}
for i, v in pairs(ARGV2) do
    local val = redis.call('hget', KEYS[1]..v, ARGV[1])
	local total = '0'
    if val then
        total = val
    end
    result[i] = {v, total}
end
return result
`)

var BalanceConsume = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 使用数量
* result 余额
*/]]
local value = tonumber(ARGV[3])
if value >= 0 then
	return redis.error_reply("number must be negative")
end
local val = redis.call('hget', KEYS[1], ARGV[1])
local total = 0
if val then
	total = tonumber(val)
end
if total + value < 0 then
	return redis.error_reply('not enough:' .. total)
end

if (redis.call('hsetnx', KEYS[2], ARGV[2], ARGV[1] .. ',' .. ARGV[3]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end

return redis.call('hincrbyfloat', KEYS[1], ARGV[1], ARGV[3])
`)

var BalanceConsumeRevoke = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 使用数量
* result  余额 
*/]]
local value = tonumber(ARGV[3])
if value >= 0 then
	return redis.error_reply("number must be negative")
end

if (redis.call('hdel', KEYS[2], ARGV[2]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end

return redis.call('hincrbyfloat', KEYS[1], ARGV[1], 0 - value)
`)

var BalanceCharge = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 充值数量
* result 余额
*/]]
local value = tonumber(ARGV[3])
if value <= 0 then
	return redis.error_reply("number must be positive")
end

if (redis.call('hsetnx', KEYS[2], ARGV[2], ARGV[1] .. ',' .. ARGV[3]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end

return redis.call('hincrbyfloat', KEYS[1], ARGV[1], ARGV[3])
`)

var BalanceChargeRevoke = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 充值数量
* result  余额 
*/]]
local value = tonumber(ARGV[3])
if value <= 0 then
	return redis.error_reply("number must be positive")
end

local val = redis.call('hget', KEYS[1], ARGV[1])
local total = 0
if val then
	total = tonumber(val)
end
if total - value < 0 then
	return redis.error_reply('not enough:' .. total)
end

if (redis.call('hdel', KEYS[2], ARGV[2]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end

return redis.call('hincrbyfloat', KEYS[1], ARGV[1], 0 - value)
`)

var BalanceConsumeIgnoreNotEnough = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 使用数量
* result 余额
*/]]
local value = tonumber(ARGV[3])
if value >= 0 then
	return redis.error_reply("number must be negative")
end
if (redis.call('hsetnx', KEYS[2], ARGV[2], ARGV[1] .. ',' .. ARGV[3]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end
return redis.call('hincrbyfloat', KEYS[1], ARGV[1], ARGV[3])
`)

var BalanceChargeRevokeIgnoreNotEnough = redis.NewScript(`
--[[/*
* KEYS[1] 余额Key
* KEYS[2] ConsumeKey
* ARGV[1] 余额编号
* ARGV[2] 流水号
* ARGV[3] 充值数量
* result  余额 
*/]]
local value = tonumber(ARGV[3])
if value <= 0 then
	return redis.error_reply("number must be positive")
end

if (redis.call('hdel', KEYS[2], ARGV[2]) == 0) then
    return redis.error_reply("exists:" ..  ARGV[2])
end

return redis.call('hincrbyfloat', KEYS[1], ARGV[1], 0 - value)
`)
