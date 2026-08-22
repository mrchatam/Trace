package main

func cmdGui(root string, args []string) int {
	return cmdLocalHTTP(root, args, localHTTPGUI)
}
