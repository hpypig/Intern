package main

import "testing"

func TestSyncProducer(t *testing.T) {
    SyncProducer("test_topic",2)
}
