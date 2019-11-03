/* Package orm (object-relational mapping) provides a set of tools on top of the KV store interface to handle
things like secondary indexes and auto-generated ID's that would otherwise need to be hand-generated on a case by
case basis.
*/
package orm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// BucketBase provides methods shared by all buckets
type BucketBase interface {
	// Has checks if key is in the bucket
	Has(ctx sdk.Context, key []byte) (bool, error)
	// PrefixScan returns an iterator between the start and end keys for the bucket
	PrefixScan(ctx sdk.Context, start []byte, end []byte, reverse bool) (Iterator, error)
	// ByIndex returns an iterator that returns objects in the bucket with the given index key
	ByIndex(ctx sdk.Context, indexName string, key []byte) (Iterator, error)
	// ByIndex returns an iterator that returns objects in the bucket with index keys between start and end. Start and
	// end can be set to nil to iterator through all values
	ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error)
}

// ExternalKeyBucket defines a bucket where the key is stored externally to the value object
type ExternalKeyBucket interface {
	BucketBase
	// GetOne deserializes the value at the given key into the pointer passed as dest
	GetOne(ctx sdk.Context, key []byte, dest interface{}) error
	// Save saves the given key value pair
	Save(ctx sdk.Context, key []byte, value interface{}) error
	// Delete deletes the value at the given key
	Delete(ctx sdk.Context, key []byte) error
}

type HasID interface {
	ID() []byte
}

// NaturalKeyBucket defines a bucket where all values implement HasID and the key is stored it the value and
// returned by the HasID method
type NaturalKeyBucket interface {
	BucketBase
	// GetOne deserializes the value into the pointer passed as dest, with the key calculated from the pointers
	// current valeu
	GetOne(ctx sdk.Context, dest HasID) error
	// Save saves the value passed in
	Save(ctx sdk.Context, value HasID) error
	// Delete deletes any value with a key corresponding the the ID of the hasID struct passed in
	Delete(ctx sdk.Context, hasID HasID) error
}

// Indexer specifies a function that takes a key value pair and returns the index key for the given index
type Indexer func(key []byte, value interface{}) (indexValue []byte, err error)

type Index struct {
	Name    string
	Indexer Indexer
}

// AutoIDBucket specifies a bucket where keys are generated via an auto-incremented interger
type AutoIDBucket interface {
	ExternalKeyBucket

	// Create auto-generates key
	Create(ctx sdk.Context, value interface{}) ([]byte, error)
}

// Iterator allows iteration through a sequence of key value pairs
type Iterator interface {
	// LoadNext loads the next value in the sequence into the pointer passed as dest and returns the key. If there
	// are no more items an error is returned
	LoadNext(dest interface{}) (key []byte, err error)
	// Release releases the iterator and should be called at the end of iteration
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
	return &indexIterator{b, ctx, it, key, key}, nil
}

func (b bucketBase) ByIndexPrefixScan(ctx sdk.Context, indexName string, start []byte, end []byte, reverse bool) (Iterator, error) {
	st := b.indexStore(ctx, indexName)
	if reverse {
		it := st.ReverseIterator(start, end)
		return &indexIterator{b, ctx,it, start, end}, nil
	} else {
		it := st.Iterator(start, end)
		return &indexIterator{b, ctx, it, start, end}, nil
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
	// TODO: delete indexes
	return nil
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

func NewAutoIDBucket(key sdk.StoreKey, bucketPrefix string, cdc *codec.Codec, indexes []Index, idGenerator func(x uint64) []byte) AutoIDBucket {
	return &autoIDBucket{externalKeyBucket{bucketBase{key, bucketPrefix, cdc, indexes}}, idGenerator}
}

type autoIDBucket struct {
	externalKeyBucket
	idGenerator func(x uint64) []byte
}

func writeUInt64(x uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, x)
	return buf[:n]
}

func readUInt64(bz []byte) (uint64, error) {
	x, n := binary.Uvarint(bz)
	if n <= 0 {
		return 0, fmt.Errorf("can't read var uint64")
	}
	return x, nil
}

func (a autoIDBucket) Create(ctx sdk.Context, value interface{}) ([]byte, error) {
	st := a.indexStore(ctx, "$")
	bz := st.Get([]byte("$"))
	var nextID uint64 = 0
	var err error
	if bz != nil {
		nextID, err = readUInt64(bz)
		if err != nil {
			return nil, err
		}
	}
	st.Set([]byte("$"), writeUInt64(nextID))
	return a.idGenerator(nextID), nil
}

type iterator struct {
	cdc *codec.Codec
	it  sdk.Iterator
}

func (i *iterator) LoadNext(dest interface{}) (key []byte, err error) {
	if !i.it.Valid() {
		return nil, fmt.Errorf("invalid")
	}
	key = i.it.Key()
	err = i.cdc.UnmarshalBinaryBare(i.it.Value(), dest)
	if err != nil {
		return nil, err
	}
	i.it.Next()
	return key, nil
}

func (i *iterator) Release() {
	i.it.Close()
}

type indexIterator struct {
	bucketBase
	ctx sdk.Context
	it    sdk.Iterator
	start []byte
	end   []byte
}

func (i indexIterator) LoadNext(dest interface{}) (key []byte, err error) {
	if !i.it.Valid() {
		return nil, fmt.Errorf("invalid")
	}
	pieces := strings.Split(string(i.it.Key()), "/")
	if len(pieces) != 2 {
		return nil, fmt.Errorf("unexpected index key")
	}
	indexPrefix, err := hex.DecodeString(pieces[0])
	if err != nil {
		return nil, err
	}
	// check out of range
	if !((i.start == nil || bytes.Compare(i.start, indexPrefix) >= 0) && (i.end == nil || bytes.Compare(indexPrefix, i.end) <= 0)) {
		return nil, fmt.Errorf("done")
	}
	key, err = hex.DecodeString(pieces[1])
	if err != nil {
		return nil, err
	}
	err = i.bucketBase.GetOne(i.ctx, key, dest)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (i indexIterator) Release() {
	i.it.Close()
}
