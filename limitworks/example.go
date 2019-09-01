package limitworks

import (
	"fmt"
	"time"
)

func testJob1() error {

	for i := 0; i < 10000; i++ {
		if i%2500 == 0 {
			fmt.Println("[job1]:", i)
			time.Sleep(time.Millisecond * 100)
		}
	}

	return nil
}

func testJob2() error {

	for i := 'a'; i < 'f'; i++ {
		fmt.Println("[job2]:", string(i))
		time.Sleep(time.Millisecond * 100)

	}
	return nil
}

func example_Concurrency1() {

	pool := NewConcurrencyEngine(1)

	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)

	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)

	time.Sleep(time.Second * 5)
}

func example_Concurrency3() {

	pool := NewConcurrencyEngine(3)

	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)

	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)
	pool.DoJob(testJob1)
	pool.DoJob(testJob2)

	time.Sleep(time.Second * 5)
}

func ConcurrencyWorkEngineExample() {
	example_Concurrency1() // job1, job2 work not parallel
	example_Concurrency3() // job1, job2 work parallel
}
