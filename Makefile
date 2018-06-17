toggle:
	@kubectl -n $(namespace) exec $(name) -- curl -XPOST http://localhost:8080/toggle -s

status:
	@kubectl -n $(namespace) exec $(name) -- curl http://localhost:8080/get -s | jq
