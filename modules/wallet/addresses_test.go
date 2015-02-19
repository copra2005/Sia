package wallet

import (
	"testing"

	"github.com/NebulousLabs/Sia/consensus"
)

// TestCoinAddress fetches a coin address from the wallet and then spends an
// output to the coin address to verify that the wallet is correctly
// recognizing coins sent to itself.
func (wt *walletTester) testCoinAddress() {
	// Get an address.
	walletAddress, _, err := wt.wallet.CoinAddress()
	if err != nil {
		wt.assistant.Tester.Fatal(err)
	}

	// Send coins to the address, in a mined block.
	siacoinInput, value := wt.assistant.FindSpendableSiacoinInput()
	txn := wt.assistant.AddSiacoinInputToTransaction(consensus.Transaction{}, siacoinInput)
	txn.SiacoinOutputs = append(txn.SiacoinOutputs, consensus.SiacoinOutput{
		Value:      value,
		UnlockHash: walletAddress,
	})
	block, err := wt.assistant.MineCurrentBlock([]consensus.Transaction{txn})
	if err != nil {
		wt.assistant.Tester.Fatal(err)
	}
	err = wt.assistant.State.AcceptBlock(block)
	if err != nil {
		wt.assistant.Tester.Fatal(err)
	}

	// Check that the wallet sees the coins.
	if wt.wallet.Balance(false).Cmp(consensus.ZeroCurrency) == 0 {
		wt.assistant.Tester.Error("wallet didn't get the coins sent to it.")
	}
}

// TestCoinAddress creates a new wallet tester and uses it to call
// testCoinAddress.
func TestCoinAddress(t *testing.T) {
	wt := newWalletTester(t)
	wt.testCoinAddress()
}
