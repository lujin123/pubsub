Package pubsub implements a simple pub-sub library.

#### install

```bash
> go get github.com/lujin123/pubsub
```

#### example

```go
package main
import "github.com/lujin123/pubsub"

func main(){
    //create a pubsub instance
    pb := pubsub.New()

    //publish message "hello world" to topic
    pb.Publish("hello world", "topic")

    //sub multi topic
    ch := pb.Sub("topic1", "topic2")

    //sub topic once
    pb.SubOnce("topic1")

    //remove sub channel from topics
    pb.Unsub(ch, "topic")

    //close pubsub
    pb.Close()
}
```