package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLRUCache(t *testing.T) {
	Convey("Given a new LRUCache", t, func() {
		cacheCap := 2
		lruCache := NewLRUCache(cacheCap)

		Convey("When setting and getting a cache element", func() {
			holder1 := &HolderInfo{HolderAddress: "0x1", IsHolder: true, BlockNumber: 1}
			lruCache.Set("0x1", holder1)
			h1, err := lruCache.Get("0x1")

			Convey("Then the element should be retrievable and match the set value", func() {
				So(err, ShouldBeNil)
				So(h1, ShouldResemble, holder1)
			})
		})

		Convey("When updating an existing cache element", func() {
			holder1 := &HolderInfo{HolderAddress: "0x1", IsHolder: true, BlockNumber: 1}
			lruCache.Set("0x1", holder1)
			holder2 := &HolderInfo{HolderAddress: "0x1", IsHolder: false, BlockNumber: 2}
			lruCache.Set("0x1", holder2)
			h2, err := lruCache.Get("0x1")

			Convey("Then the element should be updated and retrievable", func() {
				So(err, ShouldBeNil)
				So(h2, ShouldResemble, holder2)
			})
		})

		Convey("When exceeding the cache capacity", func() {
			holder1 := &HolderInfo{HolderAddress: "0x1", IsHolder: true, BlockNumber: 1}
			holder2 := &HolderInfo{HolderAddress: "0x2", IsHolder: true, BlockNumber: 1}
			holder3 := &HolderInfo{HolderAddress: "0x3", IsHolder: true, BlockNumber: 1}
			lruCache.Set("0x1", holder1)
			lruCache.Set("0x2", holder2)
			lruCache.Set("0x3", holder3)

			_, err := lruCache.Get("0x1") // izbacen

			Convey("Then the oldest element should be evicted", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Element not found")
			})
		})

		Convey("When getting all elements from the cache", func() {
			holder1 := &HolderInfo{HolderAddress: "0x1", IsHolder: true, BlockNumber: 1}
			holder2 := &HolderInfo{HolderAddress: "0x2", IsHolder: true, BlockNumber: 1}
			lruCache.Set("0x1", holder1)
			lruCache.Set("0x2", holder2)
			allElements := lruCache.GetAll()

			Convey("Then all elements should be retrievable", func() {
				So(len(allElements), ShouldEqual, 2)
				So(allElements, ShouldContain, holder1)
				So(allElements, ShouldContain, holder2)
			})
		})

		Convey("When getting all keys from the cache", func() {
			holder1 := &HolderInfo{HolderAddress: "0x1", IsHolder: true, BlockNumber: 1}
			holder2 := &HolderInfo{HolderAddress: "0x2", IsHolder: true, BlockNumber: 1}
			lruCache.Set("0x1", holder1)
			lruCache.Set("0x2", holder2)
			allKeys := lruCache.GetKeys()

			Convey("Then all keys should be retrievable", func() {
				So(len(allKeys), ShouldEqual, 2)
				So(allKeys, ShouldContain, "0x1")
				So(allKeys, ShouldContain, "0x2")
			})
		})
	})
}
