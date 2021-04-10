package mirroring

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"
)

// consider adding a rekorClient field so that the same client struct is used for all operations, potentially saving some time
type LogHandler struct {
	metadata TreeMetadata
}

// LoadFromLocal parses a JSON file to create a log handler that can be used
// to fetch and verify properties about the log.
func LoadFromLocal(filePath string) (LogHandler, error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath(filePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("configuration file could not be found.\n%s", err))
	}

	handler := LogHandler{}
	metadata, err := LoadTreeMetadata()
	if err != nil {
		return handler, err
	}
	handler.metadata = metadata
	// stored information will contain duplicates. should these duplicates be checked?
	/*
			logRoot := trilliantypes.LogRootV1{}
			logRootBytes, err := base64.StdEncoding.DecodeString(metadata.LogInfo.SignedTreeHead.LogRoot.String())
			if err != nil {
				return handler, err
			}
			err = logRoot.UnmarshalBinary(logRootBytes)
			if err != nil {
				return handler, err
			}
		logRoot.
			if logRoot.RootHash != []byte(*metadata.LogInfo.RootHash)
	*/
	return handler, nil
}

func (h *LogHandler) Save() error {
	metadata := h.metadata

	serialMetadata, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	str := viper.GetString("metadata_file_dir")
	// assumes that if file cannot be removed, it does not exist
	os.Remove(str)
	f, err := os.OpenFile(str, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(serialMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (h *LogHandler) GetPublicKey() string {
	return h.metadata.PublicKey
}

func (h *LogHandler) GetLocalTreeSize() int64 {
	return h.metadata.SavedMaxIndex + 1
}

func (h *LogHandler) GetRemoteTreeSize() (int64, error) {
	a, err := GetLogInfo()
	if err != nil {
		return 0, err
	}
	return *a.TreeSize, nil
}

func (h *LogHandler) GetRemoteRootSignature() (strfmt.Base64, error) {
	a, err := GetLogInfo()
	if err != nil {
		return nil, err
	}
	return *a.SignedTreeHead.Signature, nil
}

func (h *LogHandler) GetRootSignature() ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(h.metadata.LogInfo.SignedTreeHead.Signature.String())
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (h *LogHandler) GetRootHash() string {
	return *h.metadata.LogInfo.RootHash
}

/*
func (h *LogHandler) GetAllLeavesForType() {

}
*/

func (h *LogHandler) SetRootHash(rootHash string) {
	*h.metadata.LogInfo.RootHash = rootHash
}

func (h *LogHandler) SetPublicKey(publicKey string) {
	h.metadata.PublicKey = publicKey
}

func (h *LogHandler) SetRootSignature(sig strfmt.Base64) {
	*h.metadata.LogInfo.SignedTreeHead.Signature = sig
}

func (h *LogHandler) SetLocalTreeSize(treeSize int64) {
	h.metadata.SavedMaxIndex = treeSize - 1
}