package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func ComputeDigest(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:]))
}