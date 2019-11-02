package orm

import sdk "github.com/cosmos/cosmos-sdk/types"

type Bucket interface {
	One(ctx sdk.Context, key []byte, dest interface{})
	PrefixScan(ctx sdk.Context, prefix []byte, reverse bool) (Iterator, error)
	ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error)
	ByIndexPrefixScan(ctx sdk.Context, indexName string, prefix []byte, reverse bool) (Iterator, error)
	Save(ctx sdk.Context, key []byte, m interface{}) error
	Delete(ctx sdk.Context, key []byte) error
	Has(ctx sdk.Context, key []byte) (bool, error)
}

type AutoIDBucket interface {
	Bucket

	// Create auto-generates key
	Create(ctx sdk.Context, m interface{}) ([]byte, error)
}

type Iterator interface {

}
