# README

## Understanding Sagas in Microservices

There are different ways to orchestrate transactions across multiple services in microservices architecture, each with its own tradeoffs. Mark Richards and Neal Ford, in their great book *Software Architecture: The Hard Parts*, identified three primal forces that affect the runtime behavior and architectural characteristics of each saga:

1. **Communication**: Blocking vs Non-Blocking
2. **Coordination**: Orchestration vs Choreography
3. **Transactionality**: Atomicity vs Eventual Consistency

I attended a workshop by Mark and Neal on the O'Reilly platform. After the lecture, I wanted to get a real feeling of all these concepts and the mechanics of each saga.

## Simulating Sagas with Go

It was very easy to simulate sagas with Go, given its great support for multithreading.


### Prerequisites

- Go programming language installed
- Basic understanding of microservices architecture

### Running the Examples

Run the example simulations:

   ```sh
   go run cmd/main.go
   ```

To generate the plots

```sh
python3 plot.py
```

### Acknowledgements

- Mark Richards and Neal Ford. If it happened, see this repo; thank you so much, I learned a lot from your books and workshops. Nothing entertains me these days more than finding tradeoffs excersise.
