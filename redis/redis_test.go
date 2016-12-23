package redis

import "fmt"

func ExampleRedis() {
	client := New("192.168.0.1", 8080, "", 0)
	fmt.Println(client)

	// vals := client.HKeys("xxxxx")
	// fmt.Println(vals)

	// for _, v := range vals.Val() {
	// 	fmt.Println(v)
	// }

	// Output:
	// .
}
