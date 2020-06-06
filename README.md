# Task Ranker
Rank tasks running as docker containers in a cluster.

Task Ranker runs as a cron job on a specified schedule. Each time the task ranker is run,
it fetches data from a metric monitoring system, filters the data as required and then
submits it to a ranking strategy. The task ranking strategy uses the data received to
calibrate currently running tasks on the cluster and then rank them accordingly. The results
of the strategy are then fed back to the user through callbacks.

You will need to have a [working Golang environment running at least 1.12](https://golang.org/dl/).

### How To Use?
[Follow the instructions here](https://git-scm.com/book/en/v2/Git-Tools-Submodules) to include this repository as a git submodule.
Then follow the below configuration and run instructions.

#### Configure
Task Ranker is configured as shown below.
```go
type dummyTaskRankReceiver struct{}

func (r *dummyTaskRankReceiver) Receive(rankedTasks []entities.Task) {
	log.Println("len(rankedTasks) = ", len(rankedTasks))
}

prometheusDataFetcher, err = prometheus.NewDataFetcher(
    prometheus.WithPrometheusEndpoint("http://localhost:9090"),
    prometheus.WithFilterLabelsZeroValues([]string{"label1", "label2"}))

tRanker, err = New(
    WithDataFetcher(prometheusDataFetcher),
    WithSchedule("?/5 * * * * *"),
    WithStrategy("cpushares", &dummyTaskRankReceiver{}))
```

#### Start the Task Ranker
Once the Task Ranker has been configured, then you can start it by calling `tRanker.Start()`.

### Test
Test the module using the below command.
```commandline
go test -v ./...
```
