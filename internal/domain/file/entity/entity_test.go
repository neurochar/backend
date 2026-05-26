package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFile(t *testing.T) {
	t.Parallel()

	groupID := uuid.New()
	originalName := "resume.pdf"

	file := NewFile(groupID, "candidate_resume", true, originalName)
	assert.NotNil(t, file)

	assert.NotEqual(t, uuid.Nil, file.ID)
	assert.Equal(t, groupID, file.GroupID)
	assert.Equal(t, "candidate_resume", file.Target)
	assert.True(t, file.AssignedToTarget)
	assert.Equal(t, originalName, file.OriginalFileName)
	assert.False(t, file.UploadedToStorage)
	assert.False(t, file.ToDeleteFromStorage)
	assert.Nil(t, file.StorageFileKey)
	assert.Nil(t, file.FileMimetype)
	assert.Nil(t, file.FileHash)
}

func TestFile_SetStorageFileKey(t *testing.T) {
	t.Parallel()

	file := &File{}
	file.SetStorageFileKey("uploads/abc123.pdf")
	assert.NotNil(t, file.StorageFileKey)
	assert.Equal(t, "uploads/abc123.pdf", *file.StorageFileKey)
}

func TestFile_SetUploadedToStorage(t *testing.T) {
	t.Parallel()

	file := &File{}
	file.SetUploadedToStorage(true)
	assert.True(t, file.UploadedToStorage)
	file.SetUploadedToStorage(false)
	assert.False(t, file.UploadedToStorage)
}

func TestFile_SetAssignedToTarget(t *testing.T) {
	t.Parallel()

	file := &File{}
	file.SetAssignedToTarget(true)
	assert.True(t, file.AssignedToTarget)
	file.SetAssignedToTarget(false)
	assert.False(t, file.AssignedToTarget)
}

func TestFile_SetFileMimetype(t *testing.T) {
	t.Parallel()

	t.Run("set mimetype", func(t *testing.T) {
		file := &File{}
		mime := "application/pdf"
		file.SetFileMimetype(&mime)
		require.NotNil(t, file.FileMimetype)
		assert.Equal(t, "application/pdf", *file.FileMimetype)
	})

	t.Run("set nil", func(t *testing.T) {
		file := &File{}
		file.SetFileMimetype(nil)
		assert.Nil(t, file.FileMimetype)
	})
}

func TestFile_SetFileHash(t *testing.T) {
	t.Parallel()

	t.Run("set hash", func(t *testing.T) {
		file := &File{}
		hash := "sha256:abc"
		file.SetFileHash(&hash)
		require.NotNil(t, file.FileHash)
		assert.Equal(t, "sha256:abc", *file.FileHash)
	})

	t.Run("set nil", func(t *testing.T) {
		file := &File{}
		file.SetFileHash(nil)
		assert.Nil(t, file.FileHash)
	})
}
