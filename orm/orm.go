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

type BucketBase interface {
	PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error)
	ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error)
	ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error)
	Has(ctx sdk.Context, key []byte) (bool, error)
}

type ExternalKeyBucket interface {
	BucketBase
	GetOne(ctx sdk.Context, key []byte, dest interface{}) error
	Save(ctx sdk.Context, key []byte, m interface{}) error
	Delete(ctx sdk.Context, key []byte) error
}

type HasID interface {
	ID() []byte
}

type NaturalKeyBucket interface {
	BucketBase
	GetOne(ctx sdk.Context, dest HasID) error
	Save(ctx sdk.Context, value HasID) error
	Delete(ctx sdk.Context, hasID HasID) error
}

type Indexer func(key []byte, value interface{}) (indexValue []byte, err error)

type Index struct {
	Name    string
	Indexer Indexer
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

type bucketBase struct {
	key          sdk.StoreKey
	bucketPrefix string
	cdc          *codec.Codec
	indexes      []Index
}

func (b bucketBase) getOne(ctx sdk.Context, key []byte, dest interface{}) error {
	store := prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
	bz := store.Get(key)
	if len(bz) == 0 {
		return fmt.Errorf("not found")
	}
	return b.cdc.UnmarshalBinaryBare(bz, dest)
}

func (b bucketBase) GetOne(ctx sdk.Context, key []byte, dest interface{}) error {
	return b.getOne(ctx, key, dest)
}

func (b bucketBase) rootStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(b.bucketPrefix))
}

func (b bucketBase) PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.rootStore(ctx)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &iterator{b.cdc, it}, nil
	} else {
		it := st.Iterator(start, end)
		return &iterator{b.cdc, it}, nil
	}
}

func (b bucketBase) indexStore(ctx sdk.Context, indexName string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(b.key), []byte(fmt.Sprintf("%s/%s", b.bucketPrefix, indexName)))
}

func (b bucketBase) ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	it := st.Iterator(key, nil)
	return &indexIterator{b, it, key, key}, nil
}

func (b bucketBase) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &indexIterator{b, it, start, end}, nil
	} else {
		it := st.Iterator(start, end)
		return &indexIterator{b, it, start, end}, nil
	}
}

type externalKeyBucket struct {
	bucketBase
}

func NewExternalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) ExternalKeyBucket {
	return &externalKeyBucket{bucketBase{
		key,
		bucketPrefix,
		cdc,
		indexes,
	}}
}

func (b bucketBase) save(ctx sdk.Context, key []byte, value interface{}) error {
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
		indexStore.Set([]byte(fmt.Sprintf("%x/%x", i, key)), []byte{0})
	}
	return nil
}
func (b externalKeyBucket) Save(ctx sdk.Context, key []byte, value interface{}) error {
	return b.save(ctx, key, value)
}

func (b bucketBase) delete(ctx sdk.Context, key []byte) error {
	rootStore := b.rootStore(ctx)
	rootStore.Delete(key)
	panic("TODO: delete indexes")
}

func (b externalKeyBucket) Delete(ctx sdk.Context, key []byte) error {
	return b.delete(ctx, key)
}

func (b bucketBase) Has(ctx sdk.Context, key []byte) (bool, error) {
	rootStore := b.rootStore(ctx)
	return rootStore.Has(key), nil
}

type naturalKeyBucket struct {
	bucketBase
}

func NewNaturalKeyBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index) NaturalKeyBucket {
	return &naturalKeyBucket{bucketBase{key, bucketPrefix, cdc, indexes}}
}

func (n naturalKeyBucket) GetOne(ctx sdk.Context, dest HasID) error {
	return n.getOne(ctx, dest.ID(), dest)
}

func (n naturalKeyBucket) Save(ctx sdk.Context, value HasID) error {
	return n.save(ctx, value.ID(), value)
}

func (n naturalKeyBucket) Delete(ctx sdk.Context, hasID HasID) error {
	return n.delete(ctx, hasID.ID())
}

func NewAutoIDBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index, idGenerator  func(x uint64) []byte) AutoIDBucket {
	return &autoIDBucket{externalKeyBucket{bucketBase{key, bucketPrefix, cdc, indexes}}, idGenerator}
}

type autoIDBucket struct {
	externalKeyBucket
	idGenerator  func(x uint64) []byte
}

//func writeUInt64(x uint64) []byte {
//
//}
//
//func readUInt64(x uint64) []byte {
//
//}

func (a autoIDBucket) Create(ctx sdk.Context, value interface{}) ([]byte, error) {
	st := a.indexStore(ctx, "$")
	bz := st.Get([]byte("$"))
	var nextID uint64 = 0
	if bz != nil {
		err := a.cdc.UnmarshalBinaryBare(bz, &nextID)
		if err != nil {
			return nil, err
		}
	}
	bz, err := a.cdc.MarshalBinaryBare(nextID)
	if err != nil {
		return nil, err
	}
	st.Set([]byte("$"), bz)
	return a.idGenerator(nextID), nil
}

type iterator struct {
	cdc *codec.Codec
	it  sdk.Iterator
}

func (i *iterator) LoadNext(dest interface{}) (key []byte, err error) {
	panic("implement me")
}

func (i *iterator) Release() {
	panic("implement me")
}

type indexIterator struct {
	bucketBase
	it  sdk.Iterator
	start []byte
	end []byte
}

func (i indexIterator) LoadNext(dest interface{}) (key []byte, err error) {
	panic("implement me")
}

func (i indexIterator) Release() {
	panic("implement me")
}


