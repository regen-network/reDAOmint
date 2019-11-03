/* Package orm (object-relational mapping) provides a set of tools on top of the KV store interface to handle
things like secondary indexes and auto-generated ID's that would otherwise need to be hand-generated on a case by
case basis.
*/
package orm

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bucket interface {
	PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error)
	ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error)
	ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error)
	Has(ctx sdk.Context, key []byte) (bool, error)
}

type ExternalKeyBucket interface {
	Bucket
	GetOne(ctx sdk.Context, key []byte, dest interface{}) error
	Save(ctx sdk.Context, key []byte, m interface{}) error
	Delete(ctx sdk.Context, key []byte) error
}

type HasID interface {
	ID() []byte
}

type NaturalKeyBucket interface {
	Bucket
	GetOne(ctx sdk.Context, dest HasID) error
	Save(ctx sdk.Context, value HasID) error
	Delete(ctx sdk.Context, hasID HasID) error
}

type AutoIDBucket interface {
	ExternalKeyBucket

	// Create auto-generates key
	Create(ctx sdk.Context, value interface{}) ([]byte, error)
}

type Iterator interface {
	LoadNext(dest interface{}) (key []byte, err error)
	Release()
}

type externalKeyBucket struct {
	key          sdk.StoreKey
	bucketPrefix string
	cdc          *codec.Codec
	indexes      []Index
}

type naturalKeyBucket struct {
	key          sdk.StoreKey
	bucketPrefix string
	cdc          *codec.Codec
	indexes      []Index
}

func (n naturalKeyBucket) PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error) {
	panic("implement me")
}

func (n naturalKeyBucket) ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error) {
	panic("implement me")
}

func (n naturalKeyBucket) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	panic("implement me")
}

func (n naturalKeyBucket) Has(ctx sdk.Context, key []byte) (bool, error) {
	panic("implement me")
}

func (n naturalKeyBucket) GetOne(ctx sdk.Context, dest HasID) error {
	panic("implement me")
}

func (n naturalKeyBucket) Save(ctx sdk.Context, value HasID) error {
	panic("implement me")
}

func (n naturalKeyBucket) Delete(ctx sdk.Context, hasID HasID) error {
	panic("implement me")
}

type Indexer func(key []byte, value interface{}) (indexValue []byte, err error)

type Index struct {
	Name    string
	Indexer Indexer
}

func NewExternalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) ExternalKeyBucket {
	return &externalKeyBucket{key, bucketPrefix, cdc, indexes}
}

func NewNaturalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) NaturalKeyBucket {
	return &naturalKeyBucket{key, bucketPrefix, cdc, indexes}
}

func NewAutoIDBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index, idGenerator  func(x uint64) []byte) AutoIDBucket {
	return &autoIDBucket{key, bucketPrefix, cdc, indexes, idGenerator}
}
type autoIDBucket struct {
	key          sdk.StoreKey
	bucketPrefix string
	cdc          *codec.Codec
	indexes      []Index
	idGenerator  func(x uint64) []byte
}

func (a autoIDBucket) PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error) {
	panic("implement me")
}

func (a autoIDBucket) ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error) {
	panic("implement me")
}

func (a autoIDBucket) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	panic("implement me")
}

func (a autoIDBucket) Has(ctx sdk.Context, key []byte) (bool, error) {
	panic("implement me")
}

func (a autoIDBucket) GetOne(ctx sdk.Context, key []byte, dest interface{}) error {
	panic("implement me")
}

func (a autoIDBucket) Save(ctx sdk.Context, key []byte, m interface{}) error {
	panic("implement me")
}

func (a autoIDBucket) Delete(ctx sdk.Context, key []byte) error {
	panic("implement me")
}

func (a autoIDBucket) Create(ctx sdk.Context, value interface{}) ([]byte, error) {
	panic("implement me")
}


type iterator struct {
	cdc *codec.Codec
	it  sdk.Iterator
}

func (b externalKeyBucket) GetOne(ctx sdk.Context, key []byte, dest interface{}) error {
	store := prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
	bz := store.Get(key)
	if len(bz) == 0 {
		return fmt.Errorf("not found")
	}
	return b.cdc.UnmarshalBinaryBare(bz, dest)
}

func (b externalKeyBucket) rootStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
}

func (b externalKeyBucket) PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.rootStore(ctx)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &iterator{b.cdc, it}, nil
	} else {
		it := st.Iterator(start, end)
		return &iterator{b.cdc, it}, nil
	}
}

func (b externalKeyBucket) indexStore(ctx sdk.Context, indexName string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(fmt.Sprintf("%s/%s", b.bucketPrefix, indexName)))
}

func (b externalKeyBucket) ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	it := st.Iterator(key, key)
	return &iterator{b.cdc, it}, nil
}

func (b externalKeyBucket) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &iterator{b.cdc, it}, nil
	} else {
		it := st.Iterator(start, end)
		return &iterator{b.cdc, it}, nil
	}
}

func (b externalKeyBucket) Save(ctx sdk.Context, key []byte, value interface{}) error {
	rootStore := b.rootStore(ctx)
	bz, err := b.cdc.MarshalBinaryBare(value)
	if err != nil {
		return err
	}
	rootStore.Set(key, bz)
	for _, idx := range b.indexes {
		i, err := idx.Indexer(key, value)
		if err != nil {
			return err
		}
		indexStore := b.indexStore(ctx, idx.Name)
		indexStore.Set([]byte(fmt.Sprintf("%x, %x", i, key)), []byte{0})
	}
	return nil
}

func (b externalKeyBucket) Delete(ctx sdk.Context, key []byte) error {
	panic("TODO")
}

func (b externalKeyBucket) Has(ctx sdk.Context, key []byte) (bool, error) {
	rootStore := b.rootStore(ctx)
	return rootStore.Has(key), nil
}

func (i *iterator) LoadNext(dest interface{}) (key []byte, err error) {
	panic("implement me")
}

func (i *iterator) Release() {
	panic("implement me")
}
