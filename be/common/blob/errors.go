package blob

import "fmt"

type NotFoundError struct {
	id  string
	err error
}

func (e *NotFoundError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("blob[%s] not found: %v", e.id, e.err)
	}
	return fmt.Sprintf("blob[%s] not found", e.id)
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}

func (e *NotFoundError) Cause() error {
	return e.err
}

func NewNotFoundError(id string, err error) *NotFoundError {
	return &NotFoundError{
		id:  id,
		err: err,
	}
}

type StoreNotFoundError struct {
	store string
	err   error
}

func (e *StoreNotFoundError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("blob store[%s] not found: %v", e.store, e.err)
	}
	return fmt.Sprintf("blob store[%s] not found", e.store)
}

func (e *StoreNotFoundError) Unwrap() error {
	return e.err
}

func (e *StoreNotFoundError) Cause() error {
	return e.err
}

func NewNotStoreError(store string, err error) *StoreNotFoundError {
	return &StoreNotFoundError{
		store: store,
		err:   err,
	}
}

type BlobIdInputEmptyError struct {
}

func (*BlobIdInputEmptyError) Error() string {
	return "blob id is empty"
}

type BlobIdInputError struct {
	from string
}

func NewBlobIdInputError(from string) *BlobIdInputError {
	return &BlobIdInputError{from: from}
}

func (e *BlobIdInputError) Error() string {
	return fmt.Sprintf("from id:%s is not format: dir/id@123 regex[%s]", e.from, BlockIdRefExp.String())
}

type FolderInputError struct {
}

func (*FolderInputError) Error() string {
	return "folder name incorrect example: 'folder-1'"
}

type IdInputError struct {
}

func (*IdInputError) Error() string {
	return "ID incorrect example: 'my-file-1'"
}

type TenantNotFoundError struct {
	tenant string
}

func NewTenantNotFoundError(tenant string) *TenantNotFoundError {
	return &TenantNotFoundError{
		tenant: tenant,
	}
}

func (e TenantNotFoundError) Error() string {
	return fmt.Sprintf("tenant[%s] not found", e.tenant)
}
