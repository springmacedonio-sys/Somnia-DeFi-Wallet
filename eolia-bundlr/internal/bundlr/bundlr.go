package bundlr

import (
	"context"
	"encoding/hex"
	"eolia-bundlr/config"
	"eolia-bundlr/internal/signer"
	"eolia-bundlr/internal/types"
	"eolia-bundlr/internal/validator"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Bundlr struct {
	ChainID   *big.Int
	Signer    *signer.LocalSigner
	Queue     *OpQueue
	Validator *validator.Validator
	Ctx       context.Context
	cancel    context.CancelFunc
}

func NewBundlr(cfg *config.Config) *Bundlr {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bundlr{
		ChainID:   big.NewInt(cfg.ChainID),
		Signer:    signer.NewLocalSigner(cfg.BundlrPrivateKey, big.NewInt(cfg.ChainID)),
		Queue:     NewOpQueue(),
		Validator: validator.NewValidator(cfg.RPCURL, common.HexToAddress(cfg.EntryPoint), common.HexToAddress(cfg.BundlrAddress), common.HexToAddress(cfg.Factory)),
		Ctx:       ctx,
		cancel:    cancel,
	}
}

func (b *Bundlr) ProcessUserOperation(op *types.PackedUserOperation, opHash *common.Hash) error {
	err := b.Validator.SimulateHandleOp(op)
	if err != nil {
		return fmt.Errorf("UserOperation simulation failed: %w", err)
	}

	err = b.Queue.Add(op, opHash)
	if err != nil {
		return fmt.Errorf("OpQueue insert failed: %w", err)
	}

	b.Queue.SetAsBundled(GetOpKey(op))

	return nil
}

func (b *Bundlr) BundleAndSend() error {
	queuedOps := b.Queue.GetAll()

	if len(queuedOps) == 0 {
		return nil
	}

	var packedOps []types.PackedUserOperation
	for _, v := range queuedOps {
		if v.State != "sent" {
			packedOps = append(packedOps, *v.Op)
		}
	}

	if len(packedOps) == 0 {
		fmt.Println("No ops to bundle")
		return nil
	}

	calldata, err := b.Validator.EntryPointABI.Pack("handleOps", packedOps, b.Signer.Address())
	if err != nil {
		return fmt.Errorf("abi.Pack failed: %w", err)
	}

	nonce, err := b.Validator.Client.PendingNonceAt(context.Background(), b.Signer.Address())
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := b.Validator.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	gasLimit := uint64(8_000_000)

	tx := gtypes.NewTx(&gtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &b.Validator.EntryPoint,
		Data:     calldata,
	})

	signedTx, err := b.Signer.Sign(tx)
	if err != nil {
		return fmt.Errorf("signing tx failed: %w", err)
	}

	err = b.Validator.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("send tx failed: %w", err)
	}

	fmt.Printf("Bundled %d ops and sent tx %s\n", len(packedOps), signedTx.Hash().Hex())

	var copyQueue []QueuedOp
	for _, op := range queuedOps {
		copyQueue = append(copyQueue, *op)
	}

	for attempt := 0; attempt < 40; attempt++ {
		receipt, err := b.Validator.Client.TransactionReceipt(context.Background(), signedTx.Hash())
		if err == nil {
			eventSig := []byte("UserOperationEvent(bytes32,address,address,uint256,bool,uint256,uint256)")
			eventHash := crypto.Keccak256Hash(eventSig)
			event := b.Validator.EntryPointABI.Events["UserOperationEvent"]

			for _, log := range receipt.Logs {
				if len(log.Topics) == 0 || log.Topics[0] != eventHash {
					continue
				}

				dataMap := map[string]interface{}{}
				err := event.Inputs.UnpackIntoMap(dataMap, log.Data)
				if err != nil {
					continue
				}

				userOpHash := log.Topics[1].Hex()
				sender := log.Topics[2].Hex()
				nonce := dataMap["nonce"].(*big.Int)
				success := dataMap["success"].(bool)
				actualGasCost := dataMap["actualGasCost"].(*big.Int)
				actualGasUsed := dataMap["actualGasUsed"].(*big.Int)

				for _, queuedOp := range copyQueue {
					if queuedOp.OpHash.Hex() == userOpHash {
						fmt.Printf("TX Hash: %s, UserOpHash: %s, Sender: %s, Nonce: %s, Success: %t, ActualGasUsed: %s, ActualGasCost: %s\n",
							receipt.TxHash.Hex(), userOpHash, sender, nonce.Text(16), success, actualGasUsed.Text(16), actualGasCost.Text(16))
						b.Queue.SetAsSent(GetOpKey(queuedOp.Op), &types.UserOperationReceipt{
							UserOpHash:    userOpHash,
							Sender:        sender,
							Nonce:         "0x" + nonce.Text(16),
							Paymaster:     "0x" + hex.EncodeToString(queuedOp.Op.PaymasterAndData),
							Success:       success,
							ActualGasUsed: "0x" + actualGasUsed.Text(16),
							ActualGasCost: "0x" + actualGasCost.Text(16),
							Receipt: &types.TxReceipt{
								TransactionHash:   receipt.TxHash.Hex(),
								BlockHash:         receipt.BlockHash.Hex(),
								BlockNumber:       "0x" + receipt.BlockNumber.Text(16),
								Logs:              []any{}, // opsiyonel parse
								LogsBloom:         "0x" + hex.EncodeToString(receipt.Bloom[:]),
								GasUsed:           "0x" + strconv.FormatUint(receipt.GasUsed, 16),
								CumulativeGasUsed: "0x" + strconv.FormatUint(receipt.CumulativeGasUsed, 16),
								EffectiveGasPrice: "0x" + receipt.EffectiveGasPrice.Text(16),
							},
						})
						fmt.Printf("UserOperation %s processed successfully\n", userOpHash)
						break
					}
				}

				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (b *Bundlr) StartBundlerLoop() {
	go func() {
		for {
			err := b.BundleAndSend()
			if err != nil {
				fmt.Println("BundleAndSend error:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}
