benchmark:
	echo 'wrk.method = "POST"' > post.lua
	echo 'wrk.body = "{\"purchasedItems\":[{\"productId\":\"1c41844d-d5e9-4285-a8fc-58633fb519b9\",\"qty\":2}],\"senderContactDetail\":\"st@main.ring\",\"senderContactType\":\"email\",\"senderName\":\"string\"}"' >> post.lua
	echo 'wrk.headers["Content-Type"] = "application/json"' >> post.lua
	echo 'wrk.headers["Connection"] = "keep-alive"' >> post.lua 
	wrk -t4 -c100 -d10s -s post.lua --latency http://172.28.144.1:8080/v1/purchase