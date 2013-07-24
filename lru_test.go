package cache

import (
	"log"
	"testing"
)

func Test_LRUHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing LRUCache")
}

func Test_AddAndRetreive(t *testing.T) {
	lru := NewLRUCache(5)

	checkSize(lru, 0, t)
	lru.Add("one", 1)
	checkSize(lru, 1, t)
	val, found := lru.Get("one")

	if !found {
		t.Error("Value not found")
	}

	if val.(int) != 1 {
		t.Error("Incorrect value. Expected: 1, got: ", val)
	}

	// check for non-existent item
	val, found = lru.Get("Nothing")
	if found {
		t.Error("Found non-existant value")
	}

	if val != nil {
		t.Error("Unexpected non-nil value")
	}
}

func Test_Contains(t *testing.T) {
	lru := NewLRUCache(5)

	lru.Add("one", 1)
	lru.Add("two", 2)
	checkSize(lru, 2, t)

	if !lru.Contains("one") {
		t.Error("Conains 'one' not as epected")
	}

	if lru.Contains("three") {
		t.Error("Conains 'three' not as epected")
	}
}

func Test_AddAndDelete(t *testing.T) {
	lru := NewLRUCache(5)

	lru.Add("one", 1)
	lru.Add("two", 2)
	checkSize(lru, 2, t)

	if !lru.Delete("one") {
		t.Error("Could not delete an existing element")
	}
	checkSize(lru, 1, t)

	if lru.Delete("three") {
		t.Error("Deleting non-existent element returned true")
	}
}

func Test_Stats(t *testing.T) {
	lru := NewLRUCache(4)

	lru.Add("one", 1)
	lru.Add("two", 2)

	lru.Get("one")
	lru.Get("two")
	lru.Get("three")
	lru.Get("four")

	stats := lru.Stats()

	if stats.Capacity != 4 {
		t.Error("Capacity not as expected")
	}

	if stats.Size != 2 {
		t.Error("Size not as expected")
	}

	if stats.Usage != 0.5 {
		t.Error("Usage not as expected")
	}

	if stats.Hits != 2 {
		t.Error("Hits not as expected")
	}

	if stats.Requests != 4 {
		t.Error("Requests not as expected")
	}

	if stats.HitRate != 0.5 {
		t.Error("Hit Rate not as expected")
	}
}

func Test_Overflow(t *testing.T) {
	lru := NewLRUCache(4)

	lru.Add("one", 1)
	lru.Add("two", 2)
	lru.Add("three", 3)
	lru.Add("four", 4)

	checkSize(lru, 4, t)

	// add another, which should bump off oldest element, 'one'
	lru.Add("five", 5)

	if lru.Contains("one") {
		t.Error("Incorrect element expunged")
	}

	if !lru.Contains("five") {
		t.Error("Newest element ('five') does not exist")
	}

	// access element 'two' to make 'three' the oldest item in the list
	lru.Get("two")

	// add new element to bump 'three'
	lru.Add("six", 6)

	checkSize(lru, 4, t)

	if lru.Contains("three") {
		t.Error("Incorrect element expunged, 'three' should not exist")
	}

	if !lru.Contains("six") {
		t.Error("Newest element ('six') does not exist")
	}
}

func checkSize(lru *LRUCache, expectedSize uint, t *testing.T) {
	if lru.size != expectedSize {
		t.Error("Incorrect size. Expected: ", expectedSize, ", got: ", lru.size)
	}
}
