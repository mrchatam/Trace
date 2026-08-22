package main

func cmdServe(root string, args []string) int {
	return cmdLocalHTTP(root, args, localHTTPServe)
}
