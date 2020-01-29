package event_test

import (
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/raedahgroup/godcr-gio/event"
)

var _ = Describe("ArgumentQueue", func() {
	Context("PopString", func() {
		It("should return an error when the queue is empty", func() {
			var queue ArgumentQueue
			_, err := queue.PopString()
			Expect(err).To(Equal(ErrQueueUnderflow))
		})
		It("should return the actual strings", func() {
			err := quick.Check(func(q []string) bool {
				copied := make([]interface{}, len(q))
				for i, v := range q {
					copied[i] = v
				}
				queue := ArgumentQueue{copied}
				for _, val := range q {
					str, err := queue.PopString()
					Expect(err).To(BeNil())
					Expect(str).To(Equal(val))
				}
				return true
			}, nil)
			if err != nil {
				Fail(err.Error())
			}
		})
	})
	Context("PopInt", func() {
		It("should return an error when the queue is empty", func() {
			var queue ArgumentQueue
			_, err := queue.PopInt()
			Expect(err).To(Equal(ErrQueueUnderflow))
		})
		It("should return the actual ints", func() {
			err := quick.Check(func(q []int) bool {
				copied := make([]interface{}, len(q))
				for i, v := range q {
					copied[i] = v
				}
				queue := ArgumentQueue{copied}
				for _, val := range q {
					i, err := queue.PopInt()
					Expect(err).To(BeNil())
					Expect(i).To(Equal(val))
				}
				return true
			}, nil)
			if err != nil {
				Fail(err.Error())
			}
		})
	})
	Context("PopInt64", func() {
		It("should return an error when the queue is empty", func() {
			var queue ArgumentQueue
			_, err := queue.PopInt64()
			Expect(err).To(Equal(ErrQueueUnderflow))
		})
		It("should return the actual int64s", func() {
			err := quick.Check(func(q []int64) bool {
				copied := make([]interface{}, len(q))
				for i, v := range q {
					copied[i] = v
				}
				queue := ArgumentQueue{copied}
				for _, val := range q {
					i, err := queue.PopInt64()
					Expect(err).To(BeNil())
					Expect(i).To(Equal(val))
				}
				return true
			}, nil)
			if err != nil {
				Fail(err.Error())
			}
		})
	})
})
