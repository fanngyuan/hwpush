# hwpush
华为推送服务 Golang SDK

Production ready

```Go
	client, err := hwpush.NewClient("xxxxx", "xxxxx")

	if err != nil {
		fmt.Println(err)
	}

	message := hwpush.NewNotification("1111", "1111")
	result, err := client.SendPush(context.Background(), []string{"xxxx"}, message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

```
