package bundlr

import (
	"eolia-bundlr/internal/types"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type QueuedOp struct {
	Op        *types.PackedUserOperation
	OpHash    *common.Hash
	Timestamp time.Time
	Attempts  int
	State     string
	Receipt   *types.UserOperationReceipt
}

type OpQueue struct {
	mu  sync.Mutex
	ops map[string]*QueuedOp
}

func NewOpQueue() *OpQueue {
	return &OpQueue{
		ops: make(map[string]*QueuedOp),
	}
}

func GetOpKey(op *types.PackedUserOperation) string {
	return fmt.Sprintf("%s:%s", op.Sender, op.Nonce.String())
}

func (q *OpQueue) Add(op *types.PackedUserOperation, opHash *common.Hash) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := GetOpKey(op)
	if _, exists := q.ops[key]; exists {
		return errors.New("duplicate op")
	}

	q.ops[key] = &QueuedOp{
		Op:        op,
		Timestamp: time.Now(),
		Attempts:  0,
		State:     "pending",
		Receipt:   nil,
		OpHash:    opHash,
	}

	fmt.Printf("OP Hash: %s\n", q.ops[key].OpHash.Hex())

	return nil
}

func (q *OpQueue) GetByHash(opHash *common.Hash) (*QueuedOp, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, queuedOp := range q.ops {
		if queuedOp.OpHash.Hex() == opHash.Hex() {
			return queuedOp, nil
		}
	}
	return nil, fmt.Errorf("op not found: %s", opHash.Hex())
}

func (q *OpQueue) GetAll() []*QueuedOp {
	q.mu.Lock()
	defer q.mu.Unlock()

	var result []*QueuedOp
	for _, v := range q.ops {
		result = append(result, v)
	}
	return result
}

func (q *OpQueue) Remove(op *types.PackedUserOperation) {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := GetOpKey(op)
	delete(q.ops, key)
}

func (q *OpQueue) IncrementAttempt(op *types.PackedUserOperation) {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := GetOpKey(op)
	if queuedOp, exists := q.ops[key]; exists {
		queuedOp.Attempts++
	}
}

func (q *OpQueue) CleanupOldOps(timeout time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	for key, queuedOp := range q.ops {
		if now.Sub(queuedOp.Timestamp) > timeout {
			delete(q.ops, key)
		}
	}
}

func (q *OpQueue) Clear() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := len(q.ops)
	q.ops = make(map[string]*QueuedOp)
	return count
}

func (q *OpQueue) SetAsSent(opKey string, receipt *types.UserOperationReceipt) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if queuedOp, exists := q.ops[opKey]; exists {
		queuedOp.State = "sent"
		queuedOp.Receipt = receipt
	}
}

func (q *OpQueue) SetAsBundled(opKey string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if queuedOp, exists := q.ops[opKey]; exists {
		queuedOp.State = "bundled"
	}
}
