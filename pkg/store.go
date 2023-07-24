package giga

// A store represents a connection to a database
// with a transactional API that
type Store interface {
	Begin() (StoreTransaction, error)

	// GetInvoice returns the invoice with the given ID.
	GetInvoice(id Address) (Invoice, error)

	// ListInvoices returns a filtered list of invoices for an account.
	// pagination: next_cursor should be passed as 'cursor' on the next call (initial cursor = 0)
	// pagination: when next_cursor == 0, that is the final page of results.
	// pagination: stores CAN return < limit (or zero) items WITH next_cursor > 0 (due to filtering)
	ListInvoices(account Address, cursor int, limit int) (items []Invoice, next_cursor int, err error)

	// GetAccount returns the account with the given ForeignID.
	GetAccount(foreignID string) (Account, error)

	// GetChainState gets the last saved Best Block information (checkpoint for restart)
	// It returns giga.NotFound if the chainstate record does not exist.
	GetChainState() (ChainState, error)

	// List all unreserved UTXOs in the account's wallet.
	// Unreserved means not already being used in a pending transaction.
	GetAllUnreservedUTXOs(account Address) ([]UTXO, error)

	// Close the store.
	Close()
}

type StoreTransaction interface {
	// Commit the transaction to the store
	Commit() error
	// Rollback the transaction from the store, should
	// be a no-op of Commit has already succeeded
	Rollback() error

	// StoreInvoice stores an invoice.
	// Caller SHOULD update Account.NextExternalKey and use StoreAccount in the same StoreTransaction.
	// It returns an unspecified error if the invoice ID already exists (FIXME)
	StoreInvoice(invoice Invoice) error

	// GetInvoice returns the invoice with the given ID.
	// It returns giga.NotFound if the invoice does not exist (key: ID/address)
	GetInvoice(id Address) (Invoice, error)

	// ListInvoices returns a filtered list of invoices for an account.
	// pagination: next_cursor should be passed as 'cursor' on the next call (initial cursor = 0)
	// pagination: when next_cursor == 0, that is the final page of results.
	// pagination: stores CAN return < limit (or zero) items WITH next_cursor > 0 (due to filtering)
	ListInvoices(account Address, cursor int, limit int) (items []Invoice, next_cursor int, err error)

	// CreateAccount stores a NEW account.
	// It returns giga.AlreadyExists if the account already exists (key: ForeignID)
	CreateAccount(account Account) error

	// UpdateAccount updates an existing account.
	// It returns giga.NotFound if the account does not exist (key: ForeignID)
	// NOTE: will not update 'Privkey' or 'Address' (changes ignored or rejected)
	// NOTE: counters can only be advanced, not regressed (e.g. NextExternalKey) (ignored or rejected)
	UpdateAccount(account Account) error

	// StoreAddresses associates a list of addresses with an accountID
	StoreAddresses(accountID Address, addresses []Address, firstAddress uint32, internal bool) error

	// GetAccount returns the account with the given ForeignID.
	// It returns giga.NotFound if the account does not exist (key: ForeignID)
	GetAccount(foreignID string) (Account, error)

	// GetAccount returns the account with the given ID.
	// It returns giga.NotFound if the account does not exist (key: ID)
	GetAccountByID(ID string) (Account, error)

	// Find the accountID (HD root PKH) that owns the given Dogecoin address.
	// Also find the key index of `pkhAddress` within the HD wallet.
	FindAccountForAddress(pkhAddress Address) (accountID Address, keyIndex uint32, isInternal bool, err error)

	// List all unreserved UTXOs in the account's wallet.
	// Unreserved means not already being used in a pending transaction.
	GetAllUnreservedUTXOs(account Address) ([]UTXO, error)

	// Create an Unspent Transaction Output (at the given block height)
	CreateUTXO(txID string, vOut int64, value CoinAmount, scriptType string, pkhAddress Address, accountID Address, keyIndex uint32, isInternal bool, blockHeight int64) error

	// Mark an Unspent Transaction Output as spent (at the given block height)
	// Returns the ID of the Account that can spend this UTXO, if known to Gigawallet.
	MarkUTXOSpent(txID string, vOut int64, spentHeight int64) (accountId string, scriptAddress Address, err error)

	// What it says on the tin. We should consider
	// adding this to Store as a fast-path
	MarkInvoiceAsPaid(address Address) error

	// UpdateChainState updates the Best Block information (checkpoint for restart)
	UpdateChainState(state ChainState) error

	// RevertUTXOsAboveHeight clears chain-heights above the given height recorded in UTXOs.
	// This serves to roll back the effects of adding or spending those UTXOs.
	RevertUTXOsAboveHeight(maxValidHeight int64) error

	// RevertTxnsAboveHeight clears chain-heights above the given height recorded in Txns.
	// This serves to roll back the effects of creating or confirming those Txns.
	RevertTxnsAboveHeight(maxValidHeight int64) error

	// Increment the chain-sequence-number for multiple accounts.
	// Use this after modifying accounts' blockchain-derived state (UTXOs, TXNs)
	IncChainSeqForAccounts(accountIds []string) error

	// Find all accounts with UTXOs or TXNs created or modified above the specified block height,
	// and increment those accounts' chain-sequence-number.
	IncAccountsAffectedByRollback(maxValidHeight int64) ([]string, error)

	// Mark all UTXOs as confirmed (available to spend) after `confirmations` blocks,
	// at the current block height passed in blockHeight. This should be called each
	// time a new block is processed, i.e. blockHeight increases, but it is safe to
	// call less often (e.g. after a batch of blocks)
	ConfirmUTXOs(confirmations int, blockHeight int64) error

	// Insert (address,block-height) pairs into the Address Index.
	// The Address Index is used to find all Blocks that contain an Address.
	// Duplicates will be ignored.
	IndexAddresses(entries []AddressBlock) error
}

type ChainState struct {
	BestBlockHash   string
	BestBlockHeight int64
}

}

// Address Index: mapping from one Address to many BlockHeight (experimental)
type AddressBlock struct {
	Addr   Address
	Height int64
}
