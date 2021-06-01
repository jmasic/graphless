package clients

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/devLucian93/thesis-go/domain"
	"github.com/go-redis/redis"
)

type QueueClient interface {
	PopWorkerTasks(numberOfTasks int64) []domain.WorkerTask
	PushWorkerTasks([]int64) error
	PushWorkerTaskWeights(vertices []domain.Vertex) error
	TasksCount() int64
	SortTasks() error
	Clear()
}

func GetQueueClient(client QueueClientType) (QueueClient, error) {
	switch client {
	case QUEUE_REDIS:
		client := redis.NewClient(&redis.Options{
			Addr:     "ec2-18-216-21-233.us-east-2.compute.amazonaws.com:6380",
			Password: "",
			DB:       0,
			PoolSize: 50000,
		})
		redis := &RedisQueueClient{client}
		return redis, nil
	}

	return nil, errors.New("Unsupported queue client!")
}

const (
	QUEUE_TASK_LIST   = "list"
	QUEUE_TASK_WEIGHT = "weight"
)

var PopTasksScript = redis.NewScript(`
	local queue_len = redis.call("LLEN", KEYS[1])
	local tasks = {}
	if queue_len == 0 then
		return {}
	elseif queue_len < tonumber(ARGV[1]) then
		for i = 1, queue_len do 
			tasks[i] = redis.call("RPOP", KEYS[1])
		end
	else 
		for i = 1, tonumber(ARGV[1]) do 
			tasks[i] = redis.call("RPOP", KEYS[1])
		end
	end
	
	return tasks
`)

type RedisQueueClient struct {
	client *redis.Client
}

func (queue *RedisQueueClient) PopWorkerTasks(numberOfTasks int64) []domain.WorkerTask {

	tasksInterf, err := PopTasksScript.Run(queue.client, []string{QUEUE_TASK_LIST}, numberOfTasks).Result()

	if err != nil {
		panic(err)
	}
	tasksStrings := tasksInterf.([]interface{})
	tasks := make([]domain.WorkerTask, len(tasksStrings), len(tasksStrings))
	for i, task := range tasksStrings {
		taskString := task.(string)
		json.Unmarshal([]byte(taskString), &tasks[i])
	}
	return tasks
}

func (queue *RedisQueueClient) PushWorkerTasks(tasks []int64) error {
	//defer utils.MeasureDuration(time.Now(), "pushWorkerTasks", "")
	wtf := make([]interface{}, len(tasks))
	for i, task := range tasks {
		wtf[i] = task
	}
	return queue.client.LPush(QUEUE_TASK_LIST, wtf...).Err()
}

func (queue *RedisQueueClient) PushWorkerTaskWeights(vertices []domain.Vertex) error {
	pipe := queue.client.Pipeline()
	for _, vertex := range vertices {
		pipe.Set(fmt.Sprintf("%v:%v", QUEUE_TASK_WEIGHT, vertex.Id), len(vertex.Edges), 0)
	}
	_, err := pipe.Exec()

	return err
}

func (queue *RedisQueueClient) SortTasks() error {
	//defer utils.MeasureDuration(time.Now(), "sortTasks", "")
	return queue.client.SortStore(QUEUE_TASK_WEIGHT, QUEUE_TASK_WEIGHT, &redis.Sort{
		By:    QUEUE_TASK_WEIGHT + ":",
		Order: "DESC",
	}).Err()
}

func (queue *RedisQueueClient) TasksCount() int64 {

	count, err := queue.client.LLen(QUEUE_TASK_LIST).Result()
	if err != nil {
		panic(err)
	}
	return count
}

func (queue *RedisQueueClient) Clear() {
	err := queue.client.FlushDB().Err()
	if err != nil {
		panic(err)
	}
}
