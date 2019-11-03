package orm

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bucket interface {
	One(ctx sdk.Context, key []byte, dest interface{}) error
	PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error)
	ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error)
	ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error)
	Has(ctx sdk.Context, key []byte) (bool, error)
}

type ExternalKeyBucket interface {
	Bucket
	Save(ctx sdk.Context, key []byte, m interface{}) error
	Delete(ctx sdk.Context, key []byte) error
}

type HasID interface {
	ID() []byte
}

type NaturalKeyBucket interface {
	Bucket
	Save(ctx sdk.Context, value HasID) error
	Delete(ctx sdk.Context, hasID HasID) error
}

type AutoIDBucket interface {
	Bucket

	// Create auto-generates key
	Create(ctx sdk.Context, m interface{}) ([]byte, error)
}

type Iterator interface {
	LoadNext(dest interface{}) (key []byte, err error)
	Release()
}

type bucket struct {
	key          sdk.StoreKey
	bucketPrefix string
	cdc          *codec.Codec
	indexes      []Index
}

type Indexer func(key []byte, value interface{}) (indexValue []byte, err error)

type Index struct {
	Name    string
	Indexer Indexer
}

func NewExternalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) ExternalKeyBucket {
	panic("TODO")
}

func NewNaturalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) NaturalKeyBucket {
	panic("TODO")
}

func NewAutoIDBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) AutoIDBucket {
	panic("TODO")
}

type iterator struct {
	cdc *codec.Codec
	it  sdk.Iterator
}

func (b bucket) One(ctx sdk.Context, key []byte, dest interface{}) error {
	store := prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
	bz := store.Get(key)
	if len(bz) == 0 {
		return fmt.Errorf("not found")
	}
	return b.cdc.UnmarshalBinaryBare(bz, dest)
}

func (b bucket) rootStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
}

func (b bucket) PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.rootStore(ctx)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &iterator{b.cdc, it}, nil
	} else {
		it := st.Iterator(start, end)
		return &iterator{b.cdc, it}, nil
	}
}

func (b bucket) indexStore(ctx sdk.Context, indexName string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(fmt.Sprintf("%s/%s", b.bucketPrefix, indexName)))
}

func (b bucket) ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	it := st.Iterator(key, key)
	return &iterator{b.cdc, it}, nil
}

func (b bucket) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &iterator{b.cdc, it}, nil
	} else {
		it := st.Iterator(start, end)
		return &iterator{b.cdc, it}, nil
	}
}

func (b bucket) Save(ctx sdk.Context, key []byte, value interface{}) error {
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

func (b bucket) Delete(ctx sdk.Context, key []byte) error {
	panic("TODO")
}

func (b bucket) Has(ctx sdk.Context, key []byte) (bool, error) {
	rootStore := b.rootStore(ctx)
	return rootStore.Has(key), nil
}

func (i *iterator) LoadNext(dest interface{}) (key []byte, err error) {
	panic("implement me")
}

func (i *iterator) Release() {
	panic("implement me")
}
