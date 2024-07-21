# go-pubsub-sample

The sample about rabbit-mq, which use these libraries:
- [go/pubsub](https://pkg.go.dev/cloud.google.com/go/pubsub)
- [core-go/pubsub](https://github.com/core-go/pubsub) to wrap [go/pubsub](https://pkg.go.dev/cloud.google.com/go/pubsub)
    - Simplify the way to initialize the consumer, publisher by configurations
        - Props: when you want to change the parameter of consumer or publisher, you can change the config file, and restart Kubernetes POD, do not need to change source code and re-compile.
- [core-go/mq](https://github.com/core-go/mq) to implement this flow, which can be considered a low code tool for message queue consumer

  ![A common flow to consume a message from a message queue](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

### Similar libraries for nodejs
We also provide these libraries to support nodejs:
- [pubsub](https://github.com/core-ts/pubsub), combine with [mq-one](https://www.npmjs.com/package/mq-one) for nodejs. Example is at [pubsub-sample](https://github.com/typescript-tutorial/pubsub-sample)
