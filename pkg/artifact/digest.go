package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func ComputeDigest(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:]))
}

func ComputeDigestFromReader(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash.Sum(nil))), nil
}