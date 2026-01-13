package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	if len(cmds) == 0 {
		return
	}
	var in chan any
	out := make(chan any)
	cmds[0](in, out)
	close(out)
	in = out
	for i := 1; i < len(cmds); i++ {
		out := make(chan any)
		go cmds[i](in, out)
		in = out
	}
	if in != nil {
		for range in {
		}
	}
}

func SelectUsers(in, out chan any) {
	usersCh := make(chan User)
	var wg sync.WaitGroup
	go func() {
		for raw := range in {
			email := fmt.Sprintf("%v", raw)
			wg.Add(1)
			go func(e string) {
				defer wg.Done()
				user := GetUser(e)
				usersCh <- user
			}(email)
		}
		wg.Wait()
		close(usersCh)
	}()
	seen := make(map[uint64]struct{})
	for u := range usersCh {
		if _, ok := seen[u.ID]; !ok {
			seen[u.ID] = struct{}{}
			out <- u
		}
	}
	close(out)
}

func SelectMessages(in, out chan any) {
	batch := make([]User, 0, GetMessagesMaxUsersBatch)
	var wg sync.WaitGroup
	worker := func(usersBatch []User) {
		defer wg.Done()
		msgs, err := GetMessages(usersBatch...)
		if err != nil {
			_ = fmt.Sprintf("GetMessages error: %v", err)
			return
		}
		for _, m := range msgs {
			out <- m
		}
	}
	for raw := range in {
		u, ok := raw.(User)
		if !ok {
			continue
		}
		batch = append(batch, u)
		if len(batch) >= GetMessagesMaxUsersBatch {
			batchCopy := make([]User, len(batch))
			copy(batchCopy, batch)
			wg.Add(1)
			go worker(batchCopy)
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		batchCopy := make([]User, len(batch))
		copy(batchCopy, batch)
		wg.Add(1)
		go worker(batchCopy)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

func CheckSpam(in, out chan any) {
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	var wg sync.WaitGroup
	for raw := range in {
		id, ok := raw.(MsgID)
		if !ok {
			continue
		}
		wg.Add(1)
		go func(m MsgID) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			has, err := HasSpam(m)
			if err != nil {
				return
			}
			out <- MsgData{ID: m, HasSpam: has}
		}(id)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

func CombineResults(in, out chan any) {
	res := make([]MsgData, 0)
	for raw := range in {
		md, ok := raw.(MsgData)
		if !ok {
			continue
		}
		res = append(res, md)
	}
	sort.Slice(res, func(i int, j int) bool {

		if res[i].HasSpam != res[j].HasSpam {
			return res[i].HasSpam
		}
		return res[i].ID < res[j].ID
	})
	for _, r := range res {
		out <- fmt.Sprintf("%t %d", r.HasSpam, r.ID)
	}
	close(out)
}
